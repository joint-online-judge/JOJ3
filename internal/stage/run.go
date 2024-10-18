package stage

import (
	"fmt"
	"log/slog"
)

func Run(stages []Stage) (stageResults []StageResult, forceQuit bool, err error) {
	var executorResults []ExecutorResult
	var parserResults []ParserResult
	var tmpParserResults []ParserResult
	slog.Info("stage run start")
	for _, stage := range stages {
		slog.Info("stage start", "name", stage.Name)
		slog.Info("executor run start", "name", stage.Executor.Name)
		slog.Debug("executor run start", "name", stage.Executor.Name,
			"cmds", stage.Executor.Cmds)
		executor, ok := executorMap[stage.Executor.Name]
		if !ok {
			slog.Error("executor not found", "name", stage.Executor.Name)
			err = fmt.Errorf("executor not found: %s", stage.Executor.Name)
			return
		}
		executorResults, err = executor.Run(stage.Executor.Cmds)
		if err != nil {
			slog.Error("executor run error", "name", stage.Executor.Name, "error", err)
			return
		}
		slog.Debug("executor run done", "results", executorResults)
		for _, executorResult := range executorResults {
			slog.Debug("executor run done", "result.Files", executorResult.Files)
		}
		parserResults = []ParserResult{}
		forceQuit = false
		for _, stageParser := range stage.Parsers {
			slog.Info("parser run start", "name", stageParser.Name)
			slog.Debug("parser run start", "name", stageParser.Name,
				"conf", stageParser.Conf)
			parser, ok := parserMap[stageParser.Name]
			if !ok {
				slog.Error("parser not found", "name", stageParser.Name)
				err = fmt.Errorf("parser not found: %s", stageParser.Name)
				return
			}
			var parserForceQuit bool
			tmpParserResults, parserForceQuit, err = parser.Run(
				executorResults, stageParser.Conf)
			if err != nil {
				slog.Error("parser run error", "name", stageParser.Name, "error", err)
				return
			}
			if parserForceQuit {
				slog.Error("parser force quit", "name", stageParser.Name)
			}
			forceQuit = forceQuit || parserForceQuit
			slog.Debug("parser run done", "results", tmpParserResults)
			if len(parserResults) == 0 {
				parserResults = tmpParserResults
			} else {
				for i := range len(parserResults) {
					parserResults[i].Score += tmpParserResults[i].Score
					parserResults[i].Comment += tmpParserResults[i].Comment
				}
			}
		}
		stageResults = append(stageResults, StageResult{
			Name:      stage.Name,
			Results:   parserResults,
			ForceQuit: forceQuit,
		})
		if forceQuit {
			break
		}
	}
	return
}

func Cleanup() {
	slog.Info("stage cleanup start")
	for name, executor := range executorMap {
		err := executor.Cleanup()
		if err != nil {
			slog.Error("executor cleanup error", "name", name, "error", err)
		}
	}
}
