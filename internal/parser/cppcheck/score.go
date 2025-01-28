package cppcheck

import (
	"fmt"
	"sort"
	"strings"
)

type Severity int

const (
	ERROR Severity = iota
	WARNING
	PORTABILITY
	PERFORMANCE
	STYLE
	INFORMATION
	DEBUG
	UNKNOWN
)

func GetResult(records []Record, conf Conf) (string, int, error) {
	score := conf.Score
	comment := "### Test results summary\n\n"
	matchCount := make(map[string]int)
	scoreChange := make(map[string]int)
	for _, record := range records {
		for _, match := range conf.Matches {
			for _, keyword := range match.Keywords {
				if strings.Contains(record.Id, keyword) ||
					strings.Contains(record.Severity, keyword) {
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
	return comment, score, nil
}
