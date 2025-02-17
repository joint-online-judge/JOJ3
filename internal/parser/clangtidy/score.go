package clangtidy

import (
	"fmt"
	"sort"
	"strings"
)

func getResult(jsonMessages []JsonMessage, conf Conf) (int, string) {
	score := conf.Score
	comment := "### Test results summary\n\n"
	matchCount := make(map[string]int)
	scoreChange := make(map[string]int)
	for _, jsonMessage := range jsonMessages {
		// checkName is commas separated string here
		checkName := jsonMessage.CheckName
		for _, match := range conf.Matches {
			for _, keyword := range match.Keywords {
				if strings.Contains(checkName, keyword) {
					matchCount[keyword] += 1
					scoreChange[keyword] += -match.Score
					score += -match.Score
				}
			}
		}
	}
	type Result struct {
		Keyword     string
		Count       int
		ScoreChange int
	}
	var results []Result
	for keyword, count := range matchCount {
		results = append(results, Result{
			Keyword:     keyword,
			Count:       count,
			ScoreChange: scoreChange[keyword],
		})
	}
	sort.Slice(results, func(i, j int) bool {
		if results[i].ScoreChange != results[j].ScoreChange {
			return results[i].ScoreChange < results[j].ScoreChange
		}
		if results[i].Count != results[j].Count {
			return results[i].Count > results[j].Count
		}
		return results[i].Keyword < results[j].Keyword
	})
	for i, result := range results {
		comment += fmt.Sprintf("%d. `%s`: %d occurrence(s), %d point(s)\n",
			i+1, result.Keyword, result.Count, result.ScoreChange)
	}
	return score, comment
}
