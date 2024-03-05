package stage

import (
	"log/slog"
)

func Run(stages []Stage) []StageResult {
	var stageResults []StageResult
	for _, stage := range stages {
		slog.Debug("stage start", "name", stage.Name)
		slog.Debug("executor run start", "cmds", stage.ExecutorCmds)
		executor := executorMap[stage.ExecutorName]
		executorResults, err := executor.Run(stage.ExecutorCmds)
		if err != nil {
			slog.Error("executor run error", "name", stage.ExecutorName, "error", err)
			break
		}
		slog.Debug("executor run done", "results", executorResults)
		slog.Debug("parser run start", "config", stage.ParserConfig)
		parser := parserMap[stage.ParserName]
		parserResults, end, err := parser.Run(executorResults, stage.ParserConfig)
		if err != nil {
			slog.Error("parser run error", "name", stage.ExecutorName, "error", err)
			break
		}
		slog.Debug("parser run done", "results", parserResults)
		stageResults = append(stageResults, StageResult{
			Name:          stage.Name,
			ParserResults: parserResults,
		})
		if end {
			break
		}
	}
	return stageResults
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
