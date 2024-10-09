package clangtidy

import (
	"fmt"
	"strings"
)

func GetResult(jsonMessages []JsonMessage, conf Conf) (int, string) {
	score := conf.Score
	comment := "### Test results summary\n\n"
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
		checkName := jsonMessage.CheckName
		for _, match := range conf.Matches {
			for _, keyword := range match.Keywords {
				// TODO: The keyword in json report might also be an array, need to split it
				if strings.Contains(checkName, keyword) {
					score -= match.Score
				}
			}
		}
		listed := false
		for key := range mapping {
			if strings.Contains(checkName, key) {
				mapping[key] += 1
				listed = true
			}
		}
		if !listed {
			mapping["others"] += 1
		}
	}

	for i, key := range keys {
		comment = fmt.Sprintf("%s%d. %s: %d\n", comment, i+1, key, mapping[key])
	}
	return score, comment
}
