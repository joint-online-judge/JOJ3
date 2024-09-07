package diff

import (
	"fmt"
	"os"
	"strings"

	"focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/internal/stage"
)

type Conf struct {
	Cases []struct {
		Score      int
		StdoutPath string
	}
}

type Diff struct{}

func (*Diff) Run(results []stage.ExecutorResult, confAny any) (
	[]stage.ParserResult, bool, error,
) {
	conf, err := stage.DecodeConf[Conf](confAny)
	if err != nil {
		return nil, true, err
	}
	if len(conf.Cases) != len(results) {
		return nil, true, fmt.Errorf("cases number not match")
	}

	var res []stage.ParserResult
	for i, caseConf := range conf.Cases {
		result := results[i]
		score := 0
		var comment string

		stdout, err := os.ReadFile(caseConf.StdoutPath)
		if err != nil {
			return nil, true, err
		}

		// If no difference, assign score
		if string(stdout) == result.Files["stdout"] {
			score = caseConf.Score
			comment = fmt.Sprintf(
				"executor status: run time: %d ns, memory: %d bytes",
				result.RunTime, result.Memory,
			)
		} else {
			// Convert stdout to string and split by lines
			stdoutLines := strings.Split(string(stdout), "\n")
			resultLines := strings.Split(result.Files["stdout"], "\n")

			// Find the first difference
			diffIndex := findFirstDifferenceIndex(stdoutLines, resultLines)
			if diffIndex != -1 {
				// Get the surrounding lines from both stdout and result
				stdoutContext := getContextLines(stdoutLines, diffIndex, 10)
				resultContext := getContextLines(resultLines, diffIndex, 10)
				comment = comment + fmt.Sprintf(
					"difference found at line %d:\nExpected output:\n%s\nActual output:\n%s",
					diffIndex+1,
					resultContext,
					stdoutContext,
				)
			}
		}

		res = append(res, stage.ParserResult{
			Score:   score,
			Comment: comment,
		})
	}

	return res, false, nil
}

// findFirstDifferenceIndex finds the index of the first line where stdout and result differ.
func findFirstDifferenceIndex(stdoutLines, resultLines []string) int {
	maxLines := len(stdoutLines)
	if len(resultLines) > maxLines {
		maxLines = len(resultLines)
	}

	for i := 0; i < maxLines; i++ {
		if i >= len(stdoutLines) || i >= len(resultLines) || stdoutLines[i] != resultLines[i] {
			return i
		}
	}
	return -1
}

// getContextLines returns the surrounding lines of the specified index.
func getContextLines(lines []string, index, contextSize int) string {
	start := index - contextSize
	if start < 0 {
		start = 0
	}
	end := index + contextSize + 1
	if end > len(lines) {
		end = len(lines)
	}

	var context strings.Builder
	for i := start; i < end; i++ {
		context.WriteString(fmt.Sprintf("%d: %s\n", i+1, lines[i]))
	}
	return context.String()
}
