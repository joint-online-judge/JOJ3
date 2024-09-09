package diff

import (
	"fmt"
	"os"
	"strings"
	"unicode"

	"focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/internal/stage"
)

type Conf struct {
	Cases []struct {
		Score            int
		StdoutPath       string
		IgnoreWhitespace bool
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
		if compareChars(string(stdout), result.Files["stdout"], caseConf.IgnoreWhitespace) {
			score = caseConf.Score
		} else {
			// Convert stdout to string and split by lines
			stdoutLines := strings.Split(string(stdout), "\n")
			resultLines := strings.Split(result.Files["stdout"], "\n")

			// Find the first difference
			diffIndex := findFirstDifferenceIndex(stdoutLines, resultLines, caseConf.IgnoreWhitespace)
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

// compareChars compares two strings character by character, optionally ignoring whitespace.
func compareChars(stdout, result string, ignoreWhitespace bool) bool {
	if ignoreWhitespace {
		stdout = removeWhitespace(stdout)
		result = removeWhitespace(result)
	}
	return stdout == result
}

// removeWhitespace removes all whitespace characters from the string.
func removeWhitespace(s string) string {
	var b strings.Builder
	for _, r := range s {
		if !unicode.IsSpace(r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// findFirstDifferenceIndex finds the index of the first line where stdout and result differ.
func findFirstDifferenceIndex(stdoutLines, resultLines []string, ignoreWhitespace bool) int {
	maxLines := len(stdoutLines)
	if len(resultLines) > maxLines {
		maxLines = len(resultLines)
	}

	for i := 0; i < maxLines; i++ {
		stdoutLine := stdoutLines[i]
		resultLine := resultLines[i]

		if ignoreWhitespace {
			stdoutLine = removeWhitespace(stdoutLine)
			resultLine = removeWhitespace(resultLine)
		}

		if stdoutLine != resultLine {
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
		stdoutLine, resultLine := getLine(stdoutLines, resultLines, i)
		if stdoutLine != resultLine {
			diffBuilder.WriteString(fmt.Sprintf("- %s\n", stdoutLine))
			diffBuilder.WriteString(fmt.Sprintf("+ %s\n", resultLine))
		} else {
			diffBuilder.WriteString(fmt.Sprintf("  %s\n", stdoutLine))
		}
	}

	// Adding the diff line
	stdoutLine, resultLine := getLine(stdoutLines, resultLines, index)
	if stdoutLine != resultLine {
		diffBuilder.WriteString(fmt.Sprintf("- %s\n", stdoutLine))
		diffBuilder.WriteString(fmt.Sprintf("+ %s\n", resultLine))
	}

	// Adding context after the diff
	for i := index + 1; i < end; i++ {
		stdoutLine, resultLine := getLine(stdoutLines, resultLines, i)
		if stdoutLine != resultLine {
			diffBuilder.WriteString(fmt.Sprintf("- %s\n", stdoutLine))
			diffBuilder.WriteString(fmt.Sprintf("+ %s\n", resultLine))
		} else {
			diffBuilder.WriteString(fmt.Sprintf("  %s\n", stdoutLine))
		}
	}

	return diffBuilder.String()
}

// getLine safely retrieves lines from both stdout and result
func getLine(stdoutLines, resultLines []string, i int) (stdoutLine, resultLine string) {
	if i < len(stdoutLines) {
		stdoutLine = stdoutLines[i]
	}
	if i < len(resultLines) {
		resultLine = resultLines[i]
	}
	return
}
