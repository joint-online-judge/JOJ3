// Referenced from https://github.com/yuriisk/clang-tidy-converter/blob/master/clang_tidy_converter/parser/clang_tidy_parser.py
package clangtidy

import (
	"fmt"
	"strings"
)

type JsonMessage struct {
	Type        string                 `json:"type"`
	CheckName   string                 `json:"checkname"`
	Description string                 `json:"description"`
	Content     map[string]interface{} `json:"content"`
	Categories  []string               `json:"categories"`
	Location    map[string]interface{} `json:"location"`
	Trace       map[string]interface{} `json:"trace"`
	Severity    string                 `json:"severity"`
}

func format(messages []ClangMessage) []JsonMessage {
	formattedMessages := make([]JsonMessage, len(messages))
	for i, message := range messages {
		formattedMessages[i] = formatMessage(message)
	}
	return formattedMessages
}

func formatMessage(message ClangMessage) JsonMessage {
	result := JsonMessage{
		Type:        "issue",
		CheckName:   message.diagnosticName,
		Description: message.message,
		Content:     extractContent(message),
		Categories:  extractCategories(message),
		Location:    extractLocation(message),
		Trace:       extractTrace(message),
		Severity:    extractSeverity(message),
	}
	return result
}

func messagesToText(messages []ClangMessage) []string {
	textLines := []string{}
	for _, message := range messages {
		textLines = append(textLines, fmt.Sprintf("%s:%d:%d: %s", message.filepath, message.line, message.column, message.message))
		textLines = append(textLines, message.detailsLines...)
		textLines = append(textLines, messagesToText(message.children)...)
	}
	return textLines
}

func extractContent(message ClangMessage) map[string]interface{} {
	detailLines := ""
	for _, line := range message.detailsLines {
		if line == "" {
			continue
		}
		detailLines += (line + "\n")
	}
	for _, line := range messagesToText(message.children) {
		if line == "" {
			continue
		}
		detailLines += (line + "\n")
	}
	result := map[string]interface{}{
		"body": "```\n" + detailLines + "```",
	}
	return result
}

func removeDuplicates(list []string) []string {
	uniqueMap := make(map[string]bool)
	for _, v := range list {
		uniqueMap[v] = true
	}
	result := []string{}
	for k := range uniqueMap {
		result = append(result, k)
	}
	return result
}

func extractCategories(message ClangMessage) []string {
	bugriskCategory := "Bug Risk"
	clarityCategory := "Clarity"
	compatibilityCategory := "Compatibility"
	complexityCategory := "Complexity"
	duplicationCategory := "Duplication"
	performanceCategory := "Performance"
	securityCategory := "Security"
	styleCategory := "Style"

	categories := []string{}
	if strings.Contains(message.diagnosticName, "bugprone") {
		categories = append(categories, bugriskCategory)
	}
	if strings.Contains(message.diagnosticName, "modernize") {
		categories = append(categories, compatibilityCategory)
	}
	if strings.Contains(message.diagnosticName, "portability") {
		categories = append(categories, compatibilityCategory)
	}
	if strings.Contains(message.diagnosticName, "performance") {
		categories = append(categories, performanceCategory)
	}
	if strings.Contains(message.diagnosticName, "readability") {
		categories = append(categories, clarityCategory)
	}
	if strings.Contains(message.diagnosticName, "cloexec") {
		categories = append(categories, securityCategory)
	}
	if strings.Contains(message.diagnosticName, "security") {
		categories = append(categories, securityCategory)
	}
	if strings.Contains(message.diagnosticName, "naming") {
		categories = append(categories, styleCategory)
	}
	if strings.Contains(message.diagnosticName, "misc") {
		categories = append(categories, styleCategory)
	}
	if strings.Contains(message.diagnosticName, "cppcoreguidelines") {
		categories = append(categories, styleCategory)
	}
	if strings.Contains(message.diagnosticName, "hicpp") {
		categories = append(categories, styleCategory)
	}
	if strings.Contains(message.diagnosticName, "simplify") {
		categories = append(categories, complexityCategory)
	}
	if strings.Contains(message.diagnosticName, "redundant") {
		categories = append(categories, duplicationCategory)
	}
	if strings.HasPrefix(message.diagnosticName, "boost-use-to-string") {
		categories = append(categories, compatibilityCategory)
	}
	if len(categories) == 0 {
		categories = append(categories, bugriskCategory)
	}
	return removeDuplicates(categories)
}

func extractLocation(message ClangMessage) map[string]interface{} {
	location := map[string]interface{}{
		"path": message.filepath,
		"lines": map[string]interface{}{
			"begin": message.line,
		},
	}
	return location
}

func extractOtherLocations(message ClangMessage) []map[string]interface{} {
	locationList := []map[string]interface{}{}
	for _, child := range message.children {
		locationList = append(locationList, extractLocation(child))
		locationList = append(locationList, extractOtherLocations(child)...)
	}
	return locationList
}

func extractTrace(message ClangMessage) map[string]interface{} {
	result := map[string]interface{}{
		"locations": extractOtherLocations(message),
	}
	return result
}

func extractSeverity(message ClangMessage) string {
	switch message.level {
	case NOTE:
		return "info"
	case REMARK:
		return "minor"
	case WARNING:
		return "major"
	case ERROR:
		return "critical"
	case FATAL:
		return "blocker"
	default:
		return "unknown"
	}
}
