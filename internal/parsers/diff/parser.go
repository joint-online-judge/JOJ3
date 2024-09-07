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

		comment = fmt.Sprintf(
			"executor status: run time: %d ns, memory: %d bytes\n",
			result.RunTime, result.Memory,
		)
		// If no difference, assign score
		if string(stdout) == result.Files["stdout"] {
			score = caseConf.Score
		} else {
			// Convert stdout to string and split by lines
			stdoutLines := strings.Split(string(stdout), "\n")
			resultLines := strings.Split(result.Files["stdout"], "\n")

			// Find the first difference
			diffIndex := findFirstDifferenceIndex(stdoutLines, resultLines)
			if diffIndex != -1 {
				// Generate diff block with surrounding context
				diffOutput := generateDiffWithContext(stdoutLines, resultLines, diffIndex, 10)
				comment += fmt.Sprintf(
					"difference found at line %d:\n```diff\n%s```",
					diffIndex+1,
					diffOutput,
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

// generateDiffWithContext creates a diff block with surrounding context from stdout and result.
func generateDiffWithContext(stdoutLines, resultLines []string, index, contextSize int) string {
	var diffBuilder strings.Builder

	start := index - contextSize
	if start < 0 {
		start = 0
	}
	end := index + contextSize + 1
	if end > len(stdoutLines) {
		end = len(stdoutLines)
	}

	// Adding context before the diff
	for i := start; i < index; i++ {
		if i < len(stdoutLines) || i < len(resultLines) {
			stdoutLine := ""
			if i < len(stdoutLines) {
				stdoutLine = stdoutLines[i]
			}
			resultLine := ""
			if i < len(resultLines) {
				resultLine = resultLines[i]
			}

			if stdoutLine != resultLine {
				if stdoutLine != "" {
					diffBuilder.WriteString(fmt.Sprintf("- %s\n", stdoutLine))
				}
				if resultLine != "" {
					diffBuilder.WriteString(fmt.Sprintf("+ %s\n", resultLine))
				}
			} else {
				diffBuilder.WriteString(fmt.Sprintf("  %s\n", stdoutLine))
			}
		}
	}

	// Adding the diff line
	if index < len(stdoutLines) || index < len(resultLines) {
		stdoutLine := ""
		if index < len(stdoutLines) {
			stdoutLine = stdoutLines[index]
		}
		resultLine := ""
		if index < len(resultLines) {
			resultLine = resultLines[index]
		}

		if stdoutLine != resultLine {
			if stdoutLine != "" {
				diffBuilder.WriteString(fmt.Sprintf("- %s\n", stdoutLine))
			}
			if resultLine != "" {
				diffBuilder.WriteString(fmt.Sprintf("+ %s\n", resultLine))
			}
		}
	}

	// Adding context after the diff
	for i := index + 1; i < end; i++ {
		if i < len(stdoutLines) || i < len(resultLines) {
			stdoutLine := ""
			if i < len(stdoutLines) {
				stdoutLine = stdoutLines[i]
			}
			resultLine := ""
			if i < len(resultLines) {
				resultLine = resultLines[i]
			}

			if stdoutLine != resultLine {
				if stdoutLine != "" {
					diffBuilder.WriteString(fmt.Sprintf("- %s\n", stdoutLine))
				}
				if resultLine != "" {
					diffBuilder.WriteString(fmt.Sprintf("+ %s\n", resultLine))
				}
			} else {
				diffBuilder.WriteString(fmt.Sprintf("  %s\n", stdoutLine))
			}
		}
	}

	return diffBuilder.String()
}
