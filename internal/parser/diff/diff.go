package diff

import (
	"fmt"
	"strings"
)

// compareStrings compares two strings character by character, optionally ignoring whitespace.
func compareStrings(str1, str2 string, compareSpace bool) bool {
	if compareSpace {
		return str1 == str2
	}
	var i, j int
	l1 := len(str1)
	l2 := len(str2)
	for i < l1 && j < l2 {
		for i < l1 && isWhitespace(str1[i]) {
			i++
		}
		for j < l2 && isWhitespace(str2[j]) {
			j++
		}
		if i < l1 && j < l2 && str1[i] != str2[j] {
			return false
		}
		if i < l1 {
			i++
		}
		if j < l2 {
			j++
		}
	}
	for i < l1 && isWhitespace(str1[i]) {
		i++
	}
	for j < l2 && isWhitespace(str2[j]) {
		j++
	}
	return i == l1 && j == l2
}

func isWhitespace(b byte) bool {
	return b == ' ' ||
		b == '\t' ||
		b == '\n' ||
		b == '\r' ||
		b == '\v' ||
		b == '\f' ||
		b == 0x85 ||
		b == 0xA0
}

// myersDiff computes the Myers' diff between two slices of strings.
// src: https://github.com/cj1128/myers-diff/blob/master/main.go
func myersDiff(src, dst []string, compareSpace bool) []operation {
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
			for len(src) > t &&
				len(dst) > t &&
				compareStrings(src[t], dst[t], compareSpace) {
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

			for x < n && y < m && compareStrings(src[x], dst[y], compareSpace) {
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
	stdoutLines, resultLines []string, ops []operation, maxLength int,
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
		if maxLength > 0 && diffBuilder.Len()+len(s) > maxLength {
			remaining := maxLength - diffBuilder.Len()
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
