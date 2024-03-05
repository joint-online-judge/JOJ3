package stage

import (
	"log/slog"
)

func Run(stages []Stage) []StageResult {
	var stageResults []StageResult
	for _, stage := range stages {
		slog.Info("stage start", "name", stage.Name)
		slog.Info("executor run start", "cmds", stage.ExecutorCmds)
		executor := executorMap[stage.ExecutorName]
		executorResults, err := executor.Run(stage.ExecutorCmds)
		if err != nil {
			slog.Error("executor run error", "name", stage.ExecutorName, "error", err)
			break
		}
		slog.Info("executor run done", "results", executorResults)
		slog.Info("parser run start", "config", stage.ParserConfig)
		parser := parserMap[stage.ParserName]
		parserResults, err := parser.Run(executorResults, stage.ParserConfig)
		if err != nil {
			slog.Error("parser run error", "name", stage.ExecutorName, "error", err)
			break
		}
		slog.Info("parser run done", "results", parserResults)
		stageResults = append(stageResults, StageResult{
			Name:          stage.Name,
			ParserResults: parserResults,
		})
	}
	return stageResults
}

func Cleanup() {
	for name, executor := range executorMap {
		slog.Info("executor cleanup start", "name", name)
		err := executor.Cleanup()
		if err != nil {
			slog.Error("executor cleanup error", "name", name, "error", err)
		}
		slog.Info("executor cleanup done", "name", name)
	}
}
