package sandbox

import (
	"fmt"
	"os"
	"path/filepath"
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

func TestConvertPBCmdReturnsStandardFileReadError(t *testing.T) {
	missing := t.TempDir() + "/missing"
	_, err := convertPBCmd([]stage.Cmd{{
		Args:  []string{"true"},
		Stdin: &stage.CmdFile{Src: &missing},
	}})
	if err == nil || !strings.Contains(err.Error(), "standard file") ||
		!strings.Contains(err.Error(), "read source file") {
		t.Fatalf("convertPBCmd() error = %v, want standard source read error", err)
	}
}

func TestConvertPBCmdReturnsCopyInDirectoryWalkError(t *testing.T) {
	missing := t.TempDir() + "/missing"
	_, err := convertPBCmd([]stage.Cmd{{
		Args:      []string{"true"},
		CopyInDir: missing,
	}})
	if err == nil || !strings.Contains(err.Error(), "walk") {
		t.Fatalf("convertPBCmd() error = %v, want directory walk error", err)
	}
}

func TestConvertPBCmdPreservesMoreThanSCMMaxFDFileCount(t *testing.T) {
	const fileCount = 301
	dir := t.TempDir()
	for i := range fileCount {
		path := filepath.Join(dir, fmt.Sprintf("file-%03d", i))
		if err := os.WriteFile(path, []byte("content"), 0o600); err != nil {
			t.Fatal(err)
		}
	}

	cmds, err := convertPBCmd([]stage.Cmd{{
		Args:      []string{"true"},
		CopyIn:    make(map[string]stage.CmdFile),
		CopyInDir: dir,
	}})
	if err != nil {
		t.Fatal(err)
	}
	if len(cmds) != 1 || len(cmds[0].GetCopyIn()) != fileCount {
		t.Fatalf("converted %d commands with %d files, want 1 command with %d files",
			len(cmds), len(cmds[0].GetCopyIn()), fileCount)
	}
}
