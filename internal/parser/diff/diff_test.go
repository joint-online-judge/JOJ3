package diff

import (
	"reflect"
	"testing"
)

func TestMyersDiff(t *testing.T) {
	tests := []struct {
		name         string
		src          []string
		dst          []string
		compareSpace bool
		expected     []operation
	}{
		{
			name:         "Insert operation",
			src:          []string{"a", "b"},
			dst:          []string{"a", "b", "c"},
			compareSpace: true,
			expected:     []operation{MOVE, MOVE, INSERT},
		},
		{
			name:         "Delete operation",
			src:          []string{"a", "b", "c"},
			dst:          []string{"a", "b"},
			compareSpace: true,
			expected:     []operation{MOVE, MOVE, DELETE},
		},
		{
			name:         "No changes",
			src:          []string{"a", "b", "c"},
			dst:          []string{"a", "b", "c"},
			compareSpace: true,
			expected:     []operation{MOVE, MOVE, MOVE},
		},
		{
			name:         "Move operation",
			src:          []string{"a", "b", "c"},
			dst:          []string{"c", "a", "b"},
			compareSpace: true,
			expected:     []operation{INSERT, MOVE, MOVE, DELETE},
		},
		{
			name:         "Ignore whitespace differences",
			src:          []string{"a ", "b"},
			dst:          []string{"a", "b"},
			compareSpace: false,
			expected:     []operation{MOVE, MOVE},
		},
		{
			name:         "Consider whitespace differences",
			src:          []string{"a ", "b"},
			dst:          []string{"a", "b"},
			compareSpace: true,
			expected:     []operation{DELETE, INSERT, MOVE},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := myersDiff(test.src, test.dst, test.compareSpace)
			if !reflect.DeepEqual(result, test.expected) {
				t.Errorf("myersDiff(%v, %v, %v) = %v; want %v",
					test.src, test.dst, test.compareSpace, result, test.expected)
			}
		})
	}
}
