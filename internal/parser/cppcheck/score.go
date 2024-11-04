package cppcheck

import (
	"fmt"
	"log/slog"
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

func severityFromString(severityString string) (Severity, error) {
	switch severityString {
	case "error":
		return ERROR, nil
	case "warning":
		return WARNING, nil
	case "portability":
		return PORTABILITY, nil
	case "performance":
		return PERFORMANCE, nil
	case "style":
		return STYLE, nil
	case "information":
		return INFORMATION, nil
	case "debug":
		return DEBUG, nil
	default:
		return UNKNOWN, fmt.Errorf("unknown severity type \"%s\" for cppcheck", severityString)
	}
}

func GetResult(records []Record, conf Conf) (string, int, error) {
	score := conf.Score
	comment := "### Test results summary\n\n"
	var severityCounts [UNKNOWN + 1]int
	// TODO: remove me
	var severityScore [UNKNOWN + 1]int
	for _, match := range conf.Matches {
		severities := match.Severity
		score := match.Score
		for _, severityString := range severities {
			severity, err := severityFromString(severityString)
			if err != nil {
				return "", 0, err
			}
			severityScore[int(severity)] = score
		}
	}
	totalSeverityScore := 0
	for _, score := range severityScore {
		totalSeverityScore += score
	}
	if totalSeverityScore != 0 {
		for _, record := range records {
			if record.File == "nofile" {
				continue
			}
			severity, err := severityFromString(record.Severity)
			if err != nil {
				slog.Error("parse severity", "error", err)
			}
			severityCounts[int(severity)] += 1
			score -= severityScore[int(severity)]
		}
		comment += fmt.Sprintf("1. error: %d\n", severityCounts[0])
		comment += fmt.Sprintf("2. warning: %d\n", severityCounts[1])
		comment += fmt.Sprintf("3. portability: %d\n", severityCounts[2])
		comment += fmt.Sprintf("4. performance: %d\n", severityCounts[3])
		comment += fmt.Sprintf("5. style: %d\n", severityCounts[4])
		comment += fmt.Sprintf("6. information: %d\n", severityCounts[5])
		comment += fmt.Sprintf("7. debug: %d\n", severityCounts[6])
	}
	matchCount := make(map[string]int)
	scoreChange := make(map[string]int)
	for _, record := range records {
		for _, match := range conf.Matches {
			for _, keyword := range match.Keywords {
				if strings.Contains(record.Id, keyword) {
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
