package local

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/joint-online-judge/JOJ3/internal/stage"
)

func stringPtr(s string) *string { return &s }

func TestRunCapturesIOAndCopyOut(t *testing.T) {
	dir := t.TempDir()
	output := filepath.Join(dir, "artifact")
	results, err := (&Local{}).Run(context.Background(), []stage.Cmd{{
		Args:    []string{"/bin/sh", "-c", "read value; printf '%s:%s' \"$MARK\" \"$value\"; printf artifact > \"$1\"", "sh", output},
		Env:     []string{"MARK=env"},
		Stdin:   &stage.CmdFile{Content: stringPtr("input\n")},
		Stdout:  &stage.CmdFile{Name: stringPtr("stdout")},
		CopyOut: []string{output, filepath.Join(dir, "optional?")},
	}})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 || results[0].Status != stage.StatusAccepted {
		t.Fatalf("Run() results = %+v", results)
	}
	if got := results[0].Files["stdout"]; got != "env:input" {
		t.Fatalf("stdout = %q", got)
	}
	if got := results[0].Files[output]; got != "artifact" {
		t.Fatalf("copy-out = %q", got)
	}
}

func TestGenerateResultClassifiesErrors(t *testing.T) {
	result := (&Local{}).generateResult(errors.New("start failed"), nil, -time.Second,
		stage.Cmd{}, bytes.Buffer{}, bytes.Buffer{}, false)
	if result.Status != stage.StatusInternalError || result.ExitStatus != -1 || result.RunTime != 0 {
		t.Fatalf("generateResult() = %+v", result)
	}

	cmd := exec.Command("/bin/sh", "-c", "exit 7")
	err := cmd.Run()
	result = (&Local{}).generateResult(err, cmd.ProcessState, time.Millisecond,
		stage.Cmd{}, bytes.Buffer{}, bytes.Buffer{}, false)
	if result.Status != stage.StatusNonzeroExitStatus || result.ExitStatus != 7 {
		t.Fatalf("exit result = %+v", result)
	}
}

func TestRunReportsRequiredCopyOutError(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "missing")
	results, err := (&Local{}).Run(context.Background(), []stage.Cmd{{
		Args: []string{"/bin/true"}, CopyOut: []string{missing},
	}})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 || results[0].Status != stage.StatusFileError || results[0].Error == "" {
		t.Fatalf("Run() results = %+v", results)
	}
}

func TestRunReadsStdinFileAndRunsMultipleCommands(t *testing.T) {
	input := filepath.Join(t.TempDir(), "stdin")
	if err := os.WriteFile(input, []byte("from-file"), 0o600); err != nil {
		t.Fatal(err)
	}
	results, err := (&Local{}).Run(context.Background(), []stage.Cmd{
		{Args: []string{"/bin/true"}},
		{
			Args:   []string{"/bin/cat"},
			Stdin:  &stage.CmdFile{Src: &input},
			Stdout: &stage.CmdFile{Name: stringPtr("stdout")},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 2 || results[1].Files["stdout"] != "from-file" {
		t.Fatalf("Run() results = %+v", results)
	}
}

func TestRunSetupErrors(t *testing.T) {
	_, err := (&Local{}).Run(context.Background(), []stage.Cmd{{
		Args: []string{"/bin/cat"}, Stdin: &stage.CmdFile{Src: stringPtr(filepath.Join(t.TempDir(), "missing"))},
	}})
	if err == nil || !strings.Contains(err.Error(), "failed to open stdin file") {
		t.Fatalf("stdin error = %v", err)
	}
	_, err = (&Local{}).Run(context.Background(), []stage.Cmd{{Args: []string{filepath.Join(t.TempDir(), "missing")}}})
	if err == nil || !strings.Contains(err.Error(), "failed to start command") {
		t.Fatalf("start error = %v", err)
	}
}

func TestRunClockTimeout(t *testing.T) {
	results, err := (&Local{}).Run(context.Background(), []stage.Cmd{{
		Args: []string{"/bin/sleep", "1"}, ClockLimit: uint64(20 * time.Millisecond),
	}})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 || results[0].Status != stage.StatusTimeLimitExceeded {
		t.Fatalf("Run() results = %+v", results)
	}
}

func TestRunRejectsEmptyArgs(t *testing.T) {
	_, err := (&Local{}).Run(context.Background(), []stage.Cmd{{}})
	if err == nil || err.Error() != "command args must not be empty" {
		t.Fatalf("Run() error = %v", err)
	}
}

func TestRunCancellationKillsProcessGroup(t *testing.T) {
	dir := t.TempDir()
	ready := filepath.Join(dir, "ready")
	marker := filepath.Join(dir, "marker")
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		_, err := (&Local{}).Run(ctx, []stage.Cmd{{
			Args: []string{"/bin/sh", "-c", "touch \"$1\"; (sleep 0.3; touch \"$2\") & wait", "sh", ready, marker},
		}})
		done <- err
	}()

	deadline := time.Now().Add(2 * time.Second)
	for {
		if _, err := os.Stat(ready); err == nil {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("command did not start")
		}
		time.Sleep(10 * time.Millisecond)
	}
	cancel()
	if err := <-done; !errors.Is(err, context.Canceled) {
		t.Fatalf("Run() error = %v, want context canceled", err)
	}
	time.Sleep(500 * time.Millisecond)
	if _, err := os.Stat(marker); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("child process survived cancellation: %v", err)
	}
}
