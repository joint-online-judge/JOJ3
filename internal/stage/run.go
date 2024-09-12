package stage

import (
	"fmt"
	"log/slog"
)

func Run(stages []Stage) (stageResults []StageResult, err error) {
	var executorResults []ExecutorResult
	var parserResults []ParserResult
	var forceQuit bool
	for _, stage := range stages {
		slog.Debug("stage start", "name", stage.Name)
		slog.Debug("executor run start", "cmds", stage.ExecutorCmds)
		executor, ok := executorMap[stage.ExecutorName]
		if !ok {
			slog.Error("executor not found", "name", stage.ExecutorName)
			err = fmt.Errorf("executor not found: %s", stage.ExecutorName)
			return
		}
		executorResults, err = executor.Run(stage.ExecutorCmds)
		if err != nil {
			slog.Error("executor run error", "name", stage.ExecutorName, "error", err)
			return
		}
		slog.Debug("executor run done", "results", executorResults)
		slog.Debug("parser run start", "conf", stage.ParserConf)
		parser, ok := parserMap[stage.ParserName]
		if !ok {
			slog.Error("parser not found", "name", stage.ParserName)
			err = fmt.Errorf("parser not found: %s", stage.ParserName)
			return
		}
		parserResults, forceQuit, err = parser.Run(executorResults, stage.ParserConf)
		if err != nil {
			slog.Error("parser run error", "name", stage.ParserName, "error", err)
			return
		}
		slog.Debug("parser run done", "results", parserResults)
		stageResults = append(stageResults, StageResult{
			Name:      stage.Name,
			Results:   parserResults,
			ForceQuit: forceQuit,
		})
		if forceQuit {
			slog.Error("parser force quit", "name", stage.ParserName)
			return
		}
	}
	return
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
