// Referenced from https://github.com/yuriisk/clang-tidy-converter/blob/master/clang_tidy_converter/parser/clang_tidy_parser.py
package clangtidy

import (
	"log/slog"
	"path/filepath"
	"regexp"
	"strconv"
)

type Level int

const (
	UNKNOWN Level = iota
	NOTE
	REMARK
	WARNING
	ERROR
	FATAL
)

type ClangMessage struct {
	filepath       string
	line           int
	column         int
	level          Level
	message        string
	diagnosticName string
	detailsLines   []string
	children       []ClangMessage
}

func newClangMessage(filepath string, line int, column int, level Level, message string, diagnosticName string, detailsLines []string, children []ClangMessage) *ClangMessage {
	if detailsLines == nil {
		detailsLines = make([]string, 0)
	}
	if children == nil {
		children = make([]ClangMessage, 0)
	}

	return &ClangMessage{
		filepath:       filepath,
		line:           line,
		column:         column,
		level:          level,
		message:        message,
		diagnosticName: diagnosticName,
		detailsLines:   detailsLines,
		children:       children,
	}
}

func levelFromString(levelString string) Level {
	switch levelString {
	case "note":
		return NOTE
	case "remark":
		return REMARK
	case "warning":
		return WARNING
	case "error":
		return ERROR
	case "fatal":
		return FATAL
	default:
		return UNKNOWN
	}
}

func isIgnored(line string) bool {
	ignoreRegex := regexp.MustCompile("^error:.*$")
	return ignoreRegex.MatchString(line)
}

func parseMessage(line string) ClangMessage {
	messageRegex := regexp.MustCompile(`^(?P<filepath>.+):(?P<line>\d+):(?P<column>\d+): (?P<level>\S+): (?P<message>.*?)(?: \[(?P<diagnostic_name>.*)\])?\n$`)
	regexRes := messageRegex.FindStringSubmatch(line)
	if len(regexRes) == 0 {
		return *newClangMessage("", 0, 0, UNKNOWN, "", "", nil, nil)
	} else {
		filepath := regexRes[1]
		line, err := strconv.Atoi(regexRes[2])
		if err != nil {
			line = 0
			slog.Error("parse line", "error", err)
		}
		column, err := strconv.Atoi(regexRes[3])
		if err != nil {
			column = 0
			slog.Error("parse column", "error", err)
		}
		level := levelFromString(regexRes[4])
		message := regexRes[5]
		diagnosticName := regexRes[6]

		return ClangMessage{
			filepath:       filepath,
			line:           line,
			column:         column,
			level:          level,
			message:        message,
			diagnosticName: diagnosticName,
			detailsLines:   make([]string, 0),
			children:       make([]ClangMessage, 0),
		}
	}
}

func groupMessages(messages []ClangMessage) []ClangMessage {
	groupedMessages := make([]ClangMessage, 0)
	for _, message := range messages {
		if message.level == NOTE {
			groupedMessages[len(groupedMessages)-1].children = append(groupedMessages[len(groupedMessages)-1].children, message)
		} else {
			groupedMessages = append(groupedMessages, message)
		}
	}
	return groupedMessages
}

func convertPathsToRelative(messages *[]ClangMessage, conf Conf) {
	currentDir := conf.RootDir
	for i := range *messages {
		(*messages)[i].filepath, _ = filepath.Rel(currentDir, (*messages)[i].filepath)
	}
}

func ParseLines(lines []string, conf Conf) []ClangMessage {
	messages := make([]ClangMessage, 0)
	for _, line := range lines {
		if isIgnored(string(line)) {
			continue
		}
		message := parseMessage(string(line))
		if message.level == UNKNOWN && len(messages) > 0 {
			messages[len(messages)-1].detailsLines = append(messages[len(messages)-1].detailsLines, string(line))
		} else {
			messages = append(messages, message)
		}
	}
	convertPathsToRelative(&messages, conf)
	return groupMessages(messages)
}
