package stage

import (
	"context"
	"strings"
	"testing"
)

type contractTestExecutor struct{}

func (contractTestExecutor) Run(context.Context, []Cmd) ([]ExecutorResult, error) {
	return []ExecutorResult{{Status: StatusAccepted}}, nil
}

func (contractTestExecutor) Cleanup(context.Context) error { return nil }

type contractTestParser struct{}

func (contractTestParser) Run([]ExecutorResult, any) ([]ParserResult, bool, error) {
	return nil, false, nil
}

func TestRunRejectsParserResultCountMismatch(t *testing.T) {
	const executorName = "contract-test-executor"
	const parserName = "contract-test-parser"
	RegisterExecutor(executorName, contractTestExecutor{})
	RegisterParser(parserName, contractTestParser{})

	_, forceQuit, err := Run(context.Background(), []Stage{{
		Name:     "test",
		Executor: StageExecutor{Name: executorName, Cmds: []Cmd{{}}},
		Parsers:  []StageParser{{Name: parserName}},
	}})
	if err == nil || !strings.Contains(err.Error(), "returned 0 results for 1") {
		t.Fatalf("Run() error = %v, want result count error", err)
	}
	if forceQuit != "test" {
		t.Fatalf("Run() force quit = %q, want test", forceQuit)
	}
}
