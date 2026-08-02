package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/joint-online-judge/JOJ3/cmd/joj3/conf"
	"github.com/joint-online-judge/JOJ3/internal/stage"
)

func TestRunStagesPreservesMainAndOutputErrors(t *testing.T) {
	brokenStage := conf.ConfStage{Name: "broken"}
	brokenStage.Executor.Name = "missing-executor"
	c := &conf.Conf{
		SandboxExecServer: "localhost:5051",
		OutputPath:        t.TempDir(),
		Stages:            []conf.ConfStage{brokenStage},
	}
	_, _, err := runStages(c, nil, func([]stage.StageResult, string) {})
	if err == nil {
		t.Fatal("runStages() unexpectedly succeeded")
	}
	if !strings.Contains(err.Error(), "executor not found") ||
		!strings.Contains(err.Error(), "is a directory") {
		t.Fatalf("runStages() error = %v, want main and output errors", err)
	}
}

func TestRunStagesPreservesPhaseErrorsAndWritesMainFailure(t *testing.T) {
	broken := func(name, executor string) conf.ConfStage {
		s := conf.ConfStage{Name: name}
		s.Executor.Name = executor
		return s
	}
	outputPath := filepath.Join(t.TempDir(), "result.json")
	c := &conf.Conf{
		SandboxExecServer: "localhost:5051",
		OutputPath:        outputPath,
		PreStages:         []conf.ConfStage{broken("pre", "missing-pre")},
		Stages:            []conf.ConfStage{broken("main", "missing-main")},
		PostStages:        []conf.ConfStage{broken("post", "missing-post")},
	}
	callbackCalled := false
	_, forceQuit, err := runStages(c, nil, func(results []stage.StageResult, forceQuit string) {
		callbackCalled = true
		if len(results) != 1 || results[0].Name != "Internal Error" || forceQuit != "Internal Error" {
			t.Fatalf("callback results = %v, forceQuit = %q", results, forceQuit)
		}
	})
	if err == nil {
		t.Fatal("runStages() unexpectedly succeeded")
	}
	for _, want := range []string{"missing-pre", "missing-main", "missing-post"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("runStages() error %q does not contain %q", err, want)
		}
	}
	if forceQuit != "Internal Error" || !callbackCalled {
		t.Fatalf("forceQuit = %q, callbackCalled = %v", forceQuit, callbackCalled)
	}
	content, readErr := os.ReadFile(outputPath)
	if readErr != nil {
		t.Fatal(readErr)
	}
	if !strings.Contains(string(content), "Internal Error") ||
		!strings.Contains(string(content), "missing-main") {
		t.Fatalf("output = %s", content)
	}
}
