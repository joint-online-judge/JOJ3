package diff

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"unicode"

	"github.com/joint-online-judge/JOJ3/internal/stage"
)

// operation represents the type of edit operation.
type operation uint

const (
	INSERT operation = iota + 1
	DELETE
	MOVE
)

type Conf struct {
	PassComment string `default:"🥳Passed!\n"`
	FailComment string `default:"🧐Failed...\n"`
	Cases       []struct {
		Outputs []struct {
			Score           int
			FileName        string
			AnswerPath      string
			CompareSpace    bool
			AlwaysHide      bool
			ForceQuitOnDiff bool
			MaxDiffLength   int
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
	forceQuit := false
	for i, caseConf := range conf.Cases {
		result := results[i]
		score := 0
		comment := ""
		for _, output := range caseConf.Outputs {
			answer, err := os.ReadFile(output.AnswerPath)
			if err != nil {
				return nil, true, err
			}
			slog.Debug("compare", "filename", output.FileName,
				"answer path", output.AnswerPath,
				"actual length", len(result.Files[output.FileName]),
				"answer length", len(string(answer)))
			// If no difference, assign score
			if compareChars(string(answer), result.Files[output.FileName],
				output.CompareSpace) {
				score += output.Score
				comment += conf.PassComment
			} else {
				if output.ForceQuitOnDiff {
					forceQuit = true
				}
				comment += conf.FailComment
				comment += fmt.Sprintf("Difference found in `%s`.\n",
					output.FileName)
				if !output.AlwaysHide {
					// Convert answer to string and split by lines
					stdoutLines := strings.Split(string(answer), "\n")
					resultLines := strings.Split(
						result.Files[output.FileName], "\n")

					// Generate Myers diff
					diffOps := myersDiff(stdoutLines, resultLines)

					// Generate diff block with surrounding context
					diffOutput := generateDiffWithContext(
						stdoutLines, resultLines, diffOps, output.MaxDiffLength)
					diffOutput = strings.TrimSuffix(diffOutput, "\n  \n")
					comment += fmt.Sprintf(
						"```diff\n%s\n```\n",
						diffOutput,
					)
				} else {
					comment += "(Content hidden.)\n"
				}
			}
		}
		res = append(res, stage.ParserResult{
			Score:   score,
			Comment: comment,
		})
	}

	return res, forceQuit, nil
}

// compareChars compares two strings character by character, optionally ignoring whitespace.
func compareChars(stdout, result string, compareSpace bool) bool {
	if !compareSpace {
		stdout = removeSpace(stdout)
		result = removeSpace(result)
	}
	return stdout == result
}

// removeSpace removes all whitespace characters from the string.
func removeSpace(s string) string {
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
	for d := 0; d <= max; d += 1 {
		v := make(map[int]int, d+2)
		trace = append(trace, v)

		if d == 0 {
			t := 0
			for len(src) > t && len(dst) > t && src[t] == dst[t] {
				t += 1
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

	for d := len(trace) - 1; d > 0; d -= 1 {
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
		for i := 0; i < trace[0][0]; i += 1 {
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
func generateDiffWithContext(
	stdoutLines, resultLines []string, ops []operation, maxSize int,
) string {
	var diffBuilder strings.Builder

	srcIndex, dstIndex, lineCount := 0, 0, 0

	for _, op := range ops {
		s := ""
		switch op {
		case INSERT:
			if dstIndex < len(resultLines) {
				s = fmt.Sprintf("+ %s\n", resultLines[dstIndex])
				dstIndex += 1
			}
		case MOVE:
			if srcIndex < len(stdoutLines) {
				s = fmt.Sprintf("  %s\n", stdoutLines[srcIndex])
				srcIndex += 1
				dstIndex += 1
			}
		case DELETE:
			if srcIndex < len(stdoutLines) {
				s = fmt.Sprintf("- %s\n", stdoutLines[srcIndex])
				srcIndex += 1
				lineCount += 1
			}
		}
		if maxSize > 0 && diffBuilder.Len()+len(s) > maxSize {
			remaining := maxSize - diffBuilder.Len()
			if remaining > 0 {
				diffBuilder.WriteString(s[:remaining])
			}
			diffBuilder.WriteString("\n\n(truncated)")
			break
		}
		diffBuilder.WriteString(s)
	}

	return diffBuilder.String()
}
