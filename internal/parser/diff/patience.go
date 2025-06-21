package diff

// modified from https://github.com/peter-evans/patience

import (
	"fmt"
	"strings"
	"unicode"
)

// stringsEqual compares two strings character by character, optionally ignoring whitespace.
func stringsEqual(str1, str2 string, compareSpace bool) bool {
	if compareSpace {
		return str1 == str2
	}
	runes1 := []rune(str1)
	runes2 := []rune(str2)
	var i, j, l1, l2 int
	l1 = len(runes1)
	l2 = len(runes2)
	for i < l1 && j < l2 {
		for i < l1 && unicode.IsSpace(runes1[i]) {
			i++
		}
		for j < l2 && unicode.IsSpace(runes2[j]) {
			j++
		}
		if i >= l1 || j >= l2 {
			break
		}
		if runes1[i] != runes2[j] {
			return false
		}
		i++
		j++
	}
	for i < l1 && unicode.IsSpace(runes1[i]) {
		i++
	}
	for j < l2 && unicode.IsSpace(runes2[j]) {
		j++
	}
	return i == l1 && j == l2
}

// DiffType defines the type of a diff element.
type DiffType int8

const (
	// Delete represents a diff delete operation.
	Delete DiffType = -1
	// Insert represents a diff insert operation.
	Insert DiffType = 1
	// Equal represents no diff.
	Equal DiffType = 0
)

// DiffLine represents a single line and its diff type.
type DiffLine struct {
	Text string
	Type DiffType
}

// typeSymbol returns the associated symbol of a DiffType.
func typeSymbol(t DiffType) string {
	switch t {
	case Equal:
		return "  "
	case Insert:
		return "+ "
	case Delete:
		return "- "
	default:
		panic("unknown DiffType")
	}
}

// diffText returns the source and destination texts (all equalities, insertions and deletions).
func diffText(diffs []DiffLine) string {
	s := make([]string, len(diffs))
	for i, l := range diffs {
		s[i] = fmt.Sprintf("%s%s", typeSymbol(l.Type), l.Text)
	}
	return strings.Join(s, "\n")
}

// LCS computes the longest common subsequence of two string
// slices and returns the index pairs of the LCS.
func LCS(a, b []string, equal func(a, b string) bool) [][2]int {
	// Initialize the LCS table.
	lcs := make([][]int, len(a)+1)
	for i := 0; i <= len(a); i++ {
		lcs[i] = make([]int, len(b)+1)
	}

	// Populate the LCS table.
	for i := 1; i < len(lcs); i++ {
		for j := 1; j < len(lcs[i]); j++ {
			if equal(a[i-1], b[j-1]) {
				lcs[i][j] = lcs[i-1][j-1] + 1
			} else {
				lcs[i][j] = max(lcs[i-1][j], lcs[i][j-1])
			}
		}
	}

	// Backtrack to find the LCS.
	i, j := len(a), len(b)
	s := make([][2]int, 0, lcs[i][j])
	for i > 0 && j > 0 {
		switch {
		case equal(a[i-1], b[j-1]):
			s = append(s, [2]int{i - 1, j - 1})
			i--
			j--
		case lcs[i-1][j] > lcs[i][j-1]:
			i--
		default:
			j--
		}
	}

	// Reverse the backtracked LCS.
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}

	return s
}

// toDiffLines is a convenience function to convert a slice of strings
// to a slice of DiffLines with the specified diff type.
func toDiffLines(a []string, t DiffType) []DiffLine {
	diffs := make([]DiffLine, len(a))
	for i, l := range a {
		diffs[i] = DiffLine{l, t}
	}
	return diffs
}

// uniqueElements returns a slice of unique elements from a slice of
// strings, and a slice of the original indices of each element.
func uniqueElements(a []string) ([]string, []int) {
	m := make(map[string]int)
	for _, e := range a {
		m[e]++
	}
	elements := []string{}
	indices := []int{}
	for i, e := range a {
		if m[e] == 1 {
			elements = append(elements, e)
			indices = append(indices, i)
		}
	}
	return elements, indices
}

// patienceDiff returns the patience diff of two slices of strings.
func patienceDiff(a, b []string, equal func(a, b string) bool) []DiffLine {
	switch {
	case len(a) == 0 && len(b) == 0:
		return []DiffLine{}
	case len(a) == 0:
		return toDiffLines(b, Insert)
	case len(b) == 0:
		return toDiffLines(a, Delete)
	}

	// Find equal elements at the head of slices a and b.
	i := 0
	for i < len(a) && i < len(b) && equal(a[i], b[i]) {
		i++
	}
	if i > 0 {
		return append(
			toDiffLines(a[:i], Equal),
			patienceDiff(a[i:], b[i:], equal)...,
		)
	}

	// Find equal elements at the tail of slices a and b.
	j := 0
	for j < len(a) && j < len(b) && equal(a[len(a)-1-j], b[len(b)-1-j]) {
		j++
	}
	if j > 0 {
		return append(
			patienceDiff(a[:len(a)-j], b[:len(b)-j], equal),
			toDiffLines(a[len(a)-j:], Equal)...,
		)
	}

	// Find the longest common subsequence of unique elements in a and b.
	ua, idxa := uniqueElements(a)
	ub, idxb := uniqueElements(b)
	lcs := LCS(ua, ub, equal)

	// If the LCS is empty, the diff is all deletions and insertions.
	if len(lcs) == 0 {
		return append(toDiffLines(a, Delete), toDiffLines(b, Insert)...)
	}

	// Lookup the original indices of slices a and b.
	for i, x := range lcs {
		lcs[i][0] = idxa[x[0]]
		lcs[i][1] = idxb[x[1]]
	}

	diffs := []DiffLine{}
	ga, gb := 0, 0
	for _, ip := range lcs {
		// PatienceDiff the gaps between the lcs elements.
		diffs = append(diffs, patienceDiff(a[ga:ip[0]], b[gb:ip[1]], equal)...)
		// Append the LCS elements to the diff.
		diffs = append(diffs, DiffLine{Type: Equal, Text: a[ip[0]]})
		ga = ip[0] + 1
		gb = ip[1] + 1
	}
	// PatienceDiff the remaining elements of a and b after the final LCS element.
	diffs = append(diffs, patienceDiff(a[ga:], b[gb:], equal)...)

	return diffs
}
