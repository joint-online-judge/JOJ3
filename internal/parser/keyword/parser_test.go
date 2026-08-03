package keyword

import (
	"strings"
	"testing"

	"github.com/joint-online-judge/JOJ3/internal/stage"
)

func TestRunCountsCapsAndOrdersKeywords(t *testing.T) {
	results, forceQuit, err := (&Keyword{}).Run([]stage.ExecutorResult{{
		Files: map[string]string{"log": "error error warning"},
	}}, map[string]any{
		"score":             10,
		"files":             []any{"log"},
		"forceQuitOnDeduct": true,
		"matches": []any{
			map[string]any{"keywords": []any{"error"}, "score": 3, "maxMatchCount": 1},
			map[string]any{"keywords": []any{"warning"}, "score": 2},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 || results[0].Score != 5 || !forceQuit {
		t.Fatalf("Run() = %+v, %v", results, forceQuit)
	}
	if !strings.Contains(results[0].Comment, "`error`: 1") || !strings.Contains(results[0].Comment, "`warning`: 1") {
		t.Fatalf("comment = %q", results[0].Comment)
	}
}

func TestRunRejectsInvalidConfiguration(t *testing.T) {
	_, forceQuit, err := (&Keyword{}).Run(nil, "invalid")
	if err == nil || !forceQuit {
		t.Fatalf("Run() = forceQuit %v, error %v", forceQuit, err)
	}
}
