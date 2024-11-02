package clangtidy

import (
	"fmt"
	"strings"

	"github.com/joint-online-judge/JOJ3/pkg/utils"
)

func GetResult(jsonMessages []JsonMessage, conf Conf) (int, string) {
	score := conf.Score
	comment := "### Test results summary\n\n"
	categoryCount := map[string]int{}
	for _, jsonMessage := range jsonMessages {
		// checkName is commas separated string here
		checkName := jsonMessage.CheckName
		for _, match := range conf.Matches {
			for _, keyword := range match.Keywords {
				if strings.Contains(checkName, keyword) {
					score -= match.Score
				}
			}
		}
		checkNames := strings.Split(checkName, ",")
		for _, checkName := range checkNames {
			parts := strings.Split(checkName, "-")
			if len(parts) > 0 {
				category := parts[0]
				// checkName might be: -warnings-as-errors
				if category == "" {
					continue
				}
				categoryCount[category] += 1
			}
		}
	}
	sortedMap := utils.SortMap(categoryCount,
		func(i, j utils.Pair[string, int]) bool {
			if i.Value == j.Value {
				return i.Key < j.Key
			}
			return i.Value > j.Value
		})
	for i, kv := range sortedMap {
		comment += fmt.Sprintf("%d. %s: %d\n", i+1, kv.Key, kv.Value)
	}
	return score, comment
}
