package stage

import (
	"fmt"
	"log/slog"
)

func Run(stages []Stage) ([]StageResult, error) {
	stageResults := []StageResult{}
	for _, stage := range stages {
		slog.Debug("stage start", "name", stage.Name)
		slog.Debug("executor run start", "cmds", stage.ExecutorCmds)
		executor, ok := executorMap[stage.ExecutorName]
		if !ok {
			slog.Error("executor not found", "name", stage.ExecutorName)
			return stageResults, fmt.Errorf("executor not found: %s", stage.ExecutorName)
		}
		executorResults, err := executor.Run(stage.ExecutorCmds)
		if err != nil {
			slog.Error("executor run error", "name", stage.ExecutorName, "error", err)
			return stageResults, err
		}
		slog.Debug("executor run done", "results", executorResults)
		slog.Debug("parser run start", "conf", stage.ParserConf)
		parser, ok := parserMap[stage.ParserName]
		if !ok {
			slog.Error("parser not found", "name", stage.ParserName)
			return stageResults, err
		}
		parserResults, forceQuit, err := parser.Run(executorResults, stage.ParserConf)
		if err != nil {
			slog.Error("parser run error", "name", stage.ExecutorName, "error", err)
			return stageResults, err
		}
		slog.Debug("parser run done", "results", parserResults)
		stageResults = append(stageResults, StageResult{
			Name:      stage.Name,
			Results:   parserResults,
			ForceQuit: forceQuit,
		})
		if forceQuit {
			slog.Error("parser force quit", "name", stage.ExecutorName)
			break
		}
	}
	return stageResults, nil
}

func Cleanup() {
	for name, executor := range executorMap {
		slog.Debug("executor cleanup start", "name", name)
		err := executor.Cleanup()
		if err != nil {
			slog.Error("executor cleanup error", "name", name, "error", err)
		}
		slog.Debug("executor cleanup done", "name", name)
	}
}
