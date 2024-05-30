package cppcheck

import "fmt"

type Severity int

const (
	ERROR Severity = iota
	WARNING
	PROBABILITY
	PERFORMANCE
	STYLE
	INFORMATION
	UNKNOWN
)

func severityFromString(severityString string) (Severity, error) {
	switch severityString {
	case "error":
		return ERROR, nil
	case "warning":
		return WARNING, nil
	case "probability":
		return PROBABILITY, nil
	case "performance":
		return PERFORMANCE, nil
	case "style":
		return STYLE, nil
	case "information":
		return INFORMATION, nil
	default:
		return UNKNOWN, fmt.Errorf("unkown severity type \"%s\" for cppcheck", severityString)
	}
}

func GetResult(records []Record, conf Conf) (string, int, error) {
	result := "### Test results summary\n\n"
	var severityCounts [6]int
	var severityScore [6]int
	score := conf.Score

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

	for _, record := range records {
		severity, _ := severityFromString(record.Severity)
		severityCounts[int(severity)] += 1
		score -= severityScore[int(severity)]
	}
	result += fmt.Sprintf("1. error: %d\n", severityCounts[0])
	result += fmt.Sprintf("2. warning: %d\n", severityCounts[1])
	result += fmt.Sprintf("3. probability: %d\n", severityCounts[2])
	result += fmt.Sprintf("4. performance: %d\n", severityCounts[3])
	result += fmt.Sprintf("5. style: %d\n", severityCounts[4])
	result += fmt.Sprintf("6. information: %d", severityCounts[5])
	return result, score, nil
}
