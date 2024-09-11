package diff

import (
	"fmt"
	"os"
	"strings"
	"unicode"

	"focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/internal/stage"
)

// operation represents the type of edit operation.
type operation uint

const (
	INSERT operation = iota + 1
	DELETE
	MOVE
)

type Conf struct {
	Cases []struct {
		Outputs []struct {
			Score            int
			FileName         string
			AnswerPath       string
			IgnoreWhitespace bool
		}
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
		comment := ""
		for _, output := range caseConf.Outputs {
			answer, err := os.ReadFile(output.AnswerPath)
			if err != nil {
				return nil, true, err
			}

			// If no difference, assign score
			if compareChars(string(answer), result.Files[output.FileName], output.IgnoreWhitespace) {
				score += output.Score
			} else {
				// Convert answer to string and split by lines
				stdoutLines := strings.Split(string(answer), "\n")
				resultLines := strings.Split(result.Files[output.FileName], "\n")

				// Generate Myers diff
				diffOps := myersDiff(stdoutLines, resultLines)

				// Generate diff block with surrounding context
				diffOutput := generateDiffWithContext(stdoutLines, resultLines, diffOps)
				comment += fmt.Sprintf(
					"difference found in %s:\n```diff\n%s```\n",
					output.FileName, diffOutput,
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

// myersDiff computes the Myers' diff between two slices of strings.
// src: https://github.com/cj1128/myers-diff/blob/master/main.go
func myersDiff(src, dst []string) []operation {
	n := len(src)
	m := len(dst)
	max := n + m
	var trace []map[int]int
	var x, y int

loop:
	for d := 0; d <= max; d++ {
		v := make(map[int]int, d+2)
		trace = append(trace, v)

		if d == 0 {
			t := 0
			for len(src) > t && len(dst) > t && src[t] == dst[t] {
				t++
			}
			v[0] = t
			if t == len(src) && len(src) == len(dst) {
				break loop
			}
			continue
		}

		lastV := trace[d-1]

		for k := -d; k <= d; k += 2 {
			if k == -d || (k != d && lastV[k-1] < lastV[k+1]) {
				x = lastV[k+1]
			} else {
				x = lastV[k-1] + 1
			}

			y = x - k

			for x < n && y < m && src[x] == dst[y] {
				x, y = x+1, y+1
			}

			v[k] = x

			if x == n && y == m {
				break loop
			}
		}
	}

	var script []operation
	x = n
	y = m
	var k, prevK, prevX, prevY int

	for d := len(trace) - 1; d > 0; d-- {
		k = x - y
		lastV := trace[d-1]

		if k == -d || (k != d && lastV[k-1] < lastV[k+1]) {
			prevK = k + 1
		} else {
			prevK = k - 1
		}

		prevX = lastV[prevK]
		prevY = prevX - prevK

		for x > prevX && y > prevY {
			script = append(script, MOVE)
			x -= 1
			y -= 1
		}

		if x == prevX {
			script = append(script, INSERT)
		} else {
			script = append(script, DELETE)
		}

		x, y = prevX, prevY
	}

	if trace[0][0] != 0 {
		for i := 0; i < trace[0][0]; i++ {
			script = append(script, MOVE)
		}
	}

	return reverse(script)
}

// reverse reverses a slice of operations.
func reverse(s []operation) []operation {
	result := make([]operation, len(s))
	for i, v := range s {
		result[len(s)-1-i] = v
	}
	return result
}

// generateDiffWithContext creates a diff block with surrounding context from stdout and result.
func generateDiffWithContext(stdoutLines, resultLines []string, ops []operation) string {
	var diffBuilder strings.Builder

	srcIndex, dstIndex := 0, 0

	for _, op := range ops {
		switch op {
		case INSERT:
			if dstIndex < len(resultLines) {
				diffBuilder.WriteString(fmt.Sprintf("+ %s\n", resultLines[dstIndex]))
				dstIndex++
			}

		case MOVE:
			if srcIndex < len(stdoutLines) {
				diffBuilder.WriteString(fmt.Sprintf("  %s\n", stdoutLines[srcIndex]))
				srcIndex++
				dstIndex++
			}

		case DELETE:
			if srcIndex < len(stdoutLines) {
				diffBuilder.WriteString(fmt.Sprintf("- %s\n", stdoutLines[srcIndex]))
				srcIndex++
			}
		}
	}

	return diffBuilder.String()
}
