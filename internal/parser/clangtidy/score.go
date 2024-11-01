package clangtidy

import (
	"fmt"
	"sort"
	"strings"
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
				categoryCount[category] += 1
			}
		}
	}
	type kv struct {
		Key   string
		Value int
	}
	var ss []kv
	for k, v := range categoryCount {
		ss = append(ss, kv{k, v})
	}
	sort.Slice(ss, func(i, j int) bool {
		if ss[i].Value == ss[j].Value {
			return ss[i].Key < ss[j].Key
		}
		return ss[i].Value > ss[j].Value
	})
	for i, kv := range ss {
		comment += fmt.Sprintf("%d. %s: %d\n", i+1, kv.Key, kv.Value)
	}
	return score, comment
}
