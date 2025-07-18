package diff

import (
	"testing"
)

func TestGenerateDiffComment(t *testing.T) {
	d := &Diff{}

	t.Run("HideCommonPrefix with common prefix", func(t *testing.T) {
		answerStr := "line1\nline2\nline3\nline4"
		resultStr := "line1\nline2\nlineA\nlineB"
		output := Output{
			HideCommonPrefix: true,
		}
		comment := d.generateDiffComment(answerStr, resultStr, output)
		expected := "```diff\n(2 line(s) of common prefix hidden)\n\n- line3\n- line4\n+ lineA\n+ lineB\n```\n"
		if comment != expected {
			t.Errorf("expected %q, got %q", expected, comment)
		}
	})

	t.Run("HideCommonPrefix with no common prefix", func(t *testing.T) {
		answerStr := "line1\nline2"
		resultStr := "lineA\nlineB"
		output := Output{
			HideCommonPrefix: true,
		}
		comment := d.generateDiffComment(answerStr, resultStr, output)
		expected := "```diff\n- line1\n- line2\n+ lineA\n+ lineB\n```\n"
		if comment != expected {
			t.Errorf("expected %q, got %q", expected, comment)
		}
	})

	t.Run("HideCommonPrefix with identical content", func(t *testing.T) {
		answerStr := "line1\nline2"
		resultStr := "line1\nline2"
		output := Output{
			HideCommonPrefix: true,
		}
		comment := d.generateDiffComment(answerStr, resultStr, output)
		expected := "```diff\n(2 line(s) of common prefix hidden)\n\n\n```\n"
		if comment != expected {
			t.Errorf("expected %q, got %q", expected, comment)
		}
	})

	t.Run("HideCommonPrefix with only common prefix", func(t *testing.T) {
		answerStr := "line1\nline2"
		resultStr := "line1\nline2\nlineA"
		output := Output{
			HideCommonPrefix: true,
		}
		comment := d.generateDiffComment(answerStr, resultStr, output)
		expected := "```diff\n(2 line(s) of common prefix hidden)\n\n+ lineA\n```\n"
		if comment != expected {
			t.Errorf("expected %q, got %q", expected, comment)
		}
	})

	t.Run("MaxDiffLines truncation", func(t *testing.T) {
		answerStr := "line1\nline2\nline3\nline4\nline5\nline6\nline7\nline8"
		resultStr := "line1\nline2\nline3\nlineA\nlineB\nlineC\nlineD\nlineE"
		output := Output{
			MaxDiffLines:     3,
			HideCommonPrefix: true,
		}
		comment := d.generateDiffComment(answerStr, resultStr, output)
		expected := "```diff\n(3 line(s) of common prefix hidden)\n\n- line4\n- line5\n- line6\n+ lineA\n+ lineB\n+ lineC\n\n(truncated)\n```\n"
		if comment != expected {
			t.Errorf("expected %q, got %q", expected, comment)
		}
	})
}
