package local

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/joint-online-judge/JOJ3/internal/stage"
)

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
