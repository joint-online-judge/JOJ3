package clangtidy

import (
	"fmt"
	"strings"
)

func contains(arr []string, element string) bool {
	for i := range arr {
		// TODO: The keyword in json report might also be an array, need to split it
		if strings.Contains(arr[i], element) {
			return true
		}
	}
	return false
}

func GetScore(jsonMessages []JsonMessage, conf Conf) int {
	fullmark := conf.Score
	for _, jsonMessage := range jsonMessages {
		keyword := jsonMessage.CheckName
		for _, match := range conf.Matches {
			if contains(match.Keyword, keyword) {
				fullmark -= match.Score
				break
			}
		}
	}
	return fullmark
}

func GetComment(jsonMessages []JsonMessage) string {
	res := "### Test results summary\n\n"
	keys := [...]string{
		"codequality-unchecked-malloc-result",
		"codequality-no-global-variables",
		"codequality-no-header-guard",
		"codequality-no-fflush-stdin",
		"readability-function-size",
		"readability-duplicate-include",
		"readability-identifier-naming",
		"readability-redundant",
		"readability-misleading-indentation",
		"readability-misplaced-array-index",
		"cppcoreguidelines-init-variables",
		"bugprone-suspicious-string-compare",
		"google-global-names-in-headers",
		"clang-diagnostic",
		"clang-analyzer",
		"misc",
		"performance",
		"others",
	}
	mapping := map[string]int{}
	for _, key := range keys {
		mapping[key] = 0
	}
	for _, jsonMessage := range jsonMessages {
		keyword := jsonMessage.CheckName
		flag := true
		for key := range mapping {
			if strings.Contains(keyword, key) {
				mapping[key] += 1
				flag = false
				break
			}
		}
		if flag {
			mapping["others"] += 1
		}
	}

	for i, key := range keys {
		res = fmt.Sprintf("%s%d. %s: %d\n", res, i+1, key, mapping[key])
	}
	return res
}
