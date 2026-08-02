package sandbox

import (
	"strings"
	"testing"

	"github.com/joint-online-judge/JOJ3/internal/stage"
)

func TestConvertPBCmdReturnsSourceReadError(t *testing.T) {
	missing := t.TempDir() + "/missing"
	_, err := convertPBCmd([]stage.Cmd{{
		Args: []string{"true"},
		CopyIn: map[string]stage.CmdFile{
			"input": {Src: &missing},
		},
	}})
	if err == nil || !strings.Contains(err.Error(), "read source file") {
		t.Fatalf("convertPBCmd() error = %v, want source read error", err)
	}
}
