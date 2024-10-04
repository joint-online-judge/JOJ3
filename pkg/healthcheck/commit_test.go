package healthcheck

import (
	"testing"
)

func TestCheckMsg(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		expected bool
	}{
		{
			name:     "Valid ASCII message",
			message:  "This is a valid commit message",
			expected: true,
		},
		{
			name:     "Message with non-ASCII character",
			message:  "This message contains a non-ASCII character: é",
			expected: false,
		},
		{
			name:     "Message with ignored prefix",
			message:  "First line\nCo-authored-by: John Doe <john@example.com>\nThis is a valid message",
			expected: true,
		},
		{
			name:     "Message with ignored prefix and non-ASCII character in content",
			message:  "First line\nCo-authored-by: John Doe <john@example.com>\nThis message has a non-ASCII character: ñ",
			expected: false,
		},
		{
			name:     "Message with ignored prefix in the first line",
			message:  "Co-authored-by: Jöhn Döe <john@example.com>",
			expected: false,
		},
		{
			name:     "Multi-line message with all valid ASCII",
			message:  "First line\nSecond line\nThird line",
			expected: true,
		},
		{
			name:     "Multi-line message with non-ASCII in middle",
			message:  "First line\nSecond line with ö\nThird line",
			expected: false,
		},
		{
			name:     "Multi-line message with non-ASCII in the first line",
			message:  "First line with ö\nSecond line",
			expected: false,
		},
		{
			name:     "Message with multiple ignored prefixes",
			message:  "First line\nCo-authored-by: John Doe <john@example.com>\nReviewed-by: Jane Smith <jane@example.com>\nValid content",
			expected: true,
		},
		{
			name:     "Empty message",
			message:  "",
			expected: true,
		},
		{
			name:     "Message with only whitespace",
			message:  "   \n  \t  ",
			expected: true,
		},
		{
			name:     "Message with non-ASCII after ignored prefix",
			message:  "First line\nCo-authored-by: John Doe <john@example.com>\nReviewed-by: Jöhn Döe <john@example.com>",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := checkMsg(tt.message)
			if result != tt.expected {
				t.Errorf("checkMsg(%q) = %v, want %v", tt.message, result, tt.expected)
			}
		})
	}
}
