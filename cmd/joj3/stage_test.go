package main

import (
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
