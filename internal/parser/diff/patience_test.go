package diff

import (
	"reflect"
	"testing"
)

func TestStringsEqual(t *testing.T) {
	testCases := []struct {
		str1         string
		str2         string
		compareSpace bool
		expected     bool
	}{
		{
			str1:         "hello",
			str2:         "hello",
			compareSpace: true,
			expected:     true,
		},
		{
			str1:         "hello",
			str2:         "hello",
			compareSpace: false,
			expected:     true,
		},
		{
			str1:         "hello",
			str2:         "world",
			compareSpace: true,
			expected:     false,
		},
		{
			str1:         "hello",
			str2:         "world",
			compareSpace: false,
			expected:     false,
		},
		{
			str1:         "hello ",
			str2:         "hello",
			compareSpace: true,
			expected:     false,
		},
		{
			str1:         "hello ",
			str2:         "hello",
			compareSpace: false,
			expected:     true,
		},
		{
			str1:         "hello\t",
			str2:         "hello",
			compareSpace: true,
			expected:     false,
		},
		{
			str1:         "hello\t",
			str2:         "hello",
			compareSpace: false,
			expected:     true,
		},
		{
			str1:         "hello  world",
			str2:         "hello world",
			compareSpace: true,
			expected:     false,
		},
		{
			str1:         "hello  world",
			str2:         "hello world",
			compareSpace: false,
			expected:     true,
		},
		{
			str1:         "hello\tworld",
			str2:         "hello world",
			compareSpace: false,
			expected:     true,
		},
		{
			str1:         "hello\tworld",
			str2:         "hello world",
			compareSpace: true,
			expected:     false,
		},
		{
			str1:         "",
			str2:         "",
			compareSpace: true,
			expected:     true,
		},
		{
			str1:         "",
			str2:         "",
			compareSpace: false,
			expected:     true,
		},
		{
			str1:         " ",
			str2:         "",
			compareSpace: true,
			expected:     false,
		},
		{
			str1:         " ",
			str2:         "",
			compareSpace: false,
			expected:     true,
		},
		{
			str1:         "hello\n",
			str2:         "hello",
			compareSpace: false,
			expected:     true,
		},
		{
			str1:         "hello\n",
			str2:         "hello",
			compareSpace: true,
			expected:     false,
		},
	}

	for _, tc := range testCases {
		actual := stringsEqual(tc.str1, tc.str2, tc.compareSpace)
		if actual != tc.expected {
			t.Errorf("stringsEqual(%q, %q, %v) = %v, expected %v", tc.str1, tc.str2, tc.compareSpace, actual, tc.expected)
		}
	}
}

func TestPatienceDiff(t *testing.T) {
	equal := func(a, b string) bool {
		return a == b
	}

	testCases := []struct {
		a        []string
		b        []string
		expected []DiffLine
	}{
		{
			a:        []string{},
			b:        []string{},
			expected: []DiffLine{},
		},
		{
			a:        []string{"a", "b", "c"},
			b:        []string{"a", "b", "c"},
			expected: []DiffLine{{Text: "a", Type: Equal}, {Text: "b", Type: Equal}, {Text: "c", Type: Equal}},
		},
		{
			a:        []string{"a", "b", "c"},
			b:        []string{"a", "b", "d"},
			expected: []DiffLine{{Text: "a", Type: Equal}, {Text: "b", Type: Equal}, {Text: "c", Type: Delete}, {Text: "d", Type: Insert}},
		},
		{
			a:        []string{"a", "b", "c"},
			b:        []string{"a", "d", "c"},
			expected: []DiffLine{{Text: "a", Type: Equal}, {Text: "b", Type: Delete}, {Text: "d", Type: Insert}, {Text: "c", Type: Equal}},
		},
		{
			a:        []string{"a", "b", "c"},
			b:        []string{"d", "e", "f"},
			expected: []DiffLine{{Text: "a", Type: Delete}, {Text: "b", Type: Delete}, {Text: "c", Type: Delete}, {Text: "d", Type: Insert}, {Text: "e", Type: Insert}, {Text: "f", Type: Insert}},
		},
		{
			a:        []string{"a", "b", "c"},
			b:        []string{"a", "b", "c", "d"},
			expected: []DiffLine{{Text: "a", Type: Equal}, {Text: "b", Type: Equal}, {Text: "c", Type: Equal}, {Text: "d", Type: Insert}},
		},
		{
			a:        []string{"a", "b", "c", "d"},
			b:        []string{"a", "b", "c"},
			expected: []DiffLine{{Text: "a", Type: Equal}, {Text: "b", Type: Equal}, {Text: "c", Type: Equal}, {Text: "d", Type: Delete}},
		},
		{
			a:        []string{"a", "b", "a", "c"},
			b:        []string{"a", "b", "b", "c"},
			expected: []DiffLine{{Text: "a", Type: Equal}, {Text: "b", Type: Equal}, {Text: "a", Type: Delete}, {Text: "b", Type: Insert}, {Text: "c", Type: Equal}},
		},
		{
			a:        []string{"a", "b", "c", "b", "a"},
			b:        []string{"b", "c", "b", "a", "d"},
			expected: []DiffLine{{Text: "a", Type: Delete}, {Text: "b", Type: Equal}, {Text: "c", Type: Equal}, {Text: "b", Type: Equal}, {Text: "a", Type: Equal}, {Text: "d", Type: Insert}},
		},
	}

	for _, tc := range testCases {
		actual := patienceDiff(tc.a, tc.b, equal)
		if !reflect.DeepEqual(actual, tc.expected) {
			t.Errorf("patienceDiff(%q, %q) = %v, expected %v", tc.a, tc.b, actual, tc.expected)
		}
	}
}
