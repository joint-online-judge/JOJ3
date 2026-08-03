package stage

import (
	"context"
	"errors"
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
	originalExecutors, originalParsers := executorMap, parserMap
	executorMap, parserMap = map[string]Executor{}, map[string]Parser{}
	t.Cleanup(func() { executorMap, parserMap = originalExecutors, originalParsers })

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

type executorFunc struct {
	run     func(context.Context, []Cmd) ([]ExecutorResult, error)
	cleanup func(context.Context) error
}

func (e executorFunc) Run(ctx context.Context, cmds []Cmd) ([]ExecutorResult, error) {
	return e.run(ctx, cmds)
}

func (e executorFunc) Cleanup(ctx context.Context) error {
	if e.cleanup == nil {
		return nil
	}
	return e.cleanup(ctx)
}

type parserFunc func([]ExecutorResult, any) ([]ParserResult, bool, error)

func (p parserFunc) Run(results []ExecutorResult, conf any) ([]ParserResult, bool, error) {
	return p(results, conf)
}

func isolateRegistries(t *testing.T) {
	t.Helper()
	originalExecutors, originalParsers := executorMap, parserMap
	executorMap, parserMap = map[string]Executor{}, map[string]Parser{}
	t.Cleanup(func() { executorMap, parserMap = originalExecutors, originalParsers })
}

func TestCleanupJoinsExecutorErrors(t *testing.T) {
	isolateRegistries(t)
	errOne := errors.New("cleanup one")
	errTwo := errors.New("cleanup two")
	RegisterExecutor("one", executorFunc{
		run:     func(context.Context, []Cmd) ([]ExecutorResult, error) { return nil, nil },
		cleanup: func(context.Context) error { return errOne },
	})
	RegisterExecutor("two", executorFunc{
		run:     func(context.Context, []Cmd) ([]ExecutorResult, error) { return nil, nil },
		cleanup: func(context.Context) error { return errTwo },
	})

	err := Cleanup(context.Background())
	if !errors.Is(err, errOne) || !errors.Is(err, errTwo) {
		t.Fatalf("Cleanup() error = %v, want both cleanup errors", err)
	}
}

func TestRunStopsAfterParserContractViolation(t *testing.T) {
	isolateRegistries(t)
	RegisterExecutor("executor", executorFunc{
		run: func(context.Context, []Cmd) ([]ExecutorResult, error) {
			return []ExecutorResult{{Status: StatusAccepted}}, nil
		},
	})
	RegisterParser("bad", parserFunc(func([]ExecutorResult, any) ([]ParserResult, bool, error) {
		return nil, false, nil
	}))
	secondCalled := false
	RegisterParser("second", parserFunc(func([]ExecutorResult, any) ([]ParserResult, bool, error) {
		secondCalled = true
		return []ParserResult{{}}, false, nil
	}))

	_, forceQuit, err := Run(context.Background(), []Stage{{
		Name:     "contract",
		Executor: StageExecutor{Name: "executor", Cmds: []Cmd{{}}},
		Parsers:  []StageParser{{Name: "bad"}, {Name: "second"}},
	}})
	if err == nil || forceQuit != "contract" {
		t.Fatalf("Run() = forceQuit %q, error %v", forceQuit, err)
	}
	if secondCalled {
		t.Fatal("parser after contract violation was called")
	}
}

func TestRunSetsForceQuitForMissingComponents(t *testing.T) {
	tests := []struct {
		name  string
		stage Stage
	}{
		{name: "executor", stage: Stage{Name: "missing-executor", Executor: StageExecutor{Name: "unknown"}}},
		{
			name: "parser",
			stage: Stage{
				Name:     "missing-parser",
				Executor: StageExecutor{Name: "executor", Cmds: []Cmd{{}}},
				Parsers:  []StageParser{{Name: "unknown"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isolateRegistries(t)
			RegisterExecutor("executor", contractTestExecutor{})
			_, forceQuit, err := Run(context.Background(), []Stage{tt.stage})
			if err == nil || forceQuit != tt.stage.Name {
				t.Fatalf("Run() = forceQuit %q, error %v", forceQuit, err)
			}
		})
	}
}
