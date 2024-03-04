package stage

import (
	"log/slog"
	"os"

	"github.com/pelletier/go-toml/v2"
)

func ParseStages(tomlConfig string) []Stage {
	var stagesConfig StagesConfig
	err := toml.Unmarshal([]byte(tomlConfig), &stagesConfig)
	if err != nil {
		slog.Error("parse stages config", "error", err)
		os.Exit(1)
	}
	stages := []Stage{}
	for _, stage := range stagesConfig.Stages {
		stages = append(stages, Stage{
			Name:         stage.Name,
			ExecutorName: stage.Executor.Name,
			Executor:     executorMap[stage.Executor.Name],
			ExecutorCmd:  stage.Executor.With,
			ParserName:   stage.Parser.Name,
			Parser:       parserMap[stage.Parser.Name],
			ParserConfig: stage.Parser.With,
		})
	}
	return stages
}

func Run(stages []Stage) []StageResult {
	var parserResults []StageResult
	for _, stage := range stages {
		slog.Info("stage start", "name", stage.Name)
		executorResult, err := stage.Executor.Run(stage.ExecutorCmd)
		if err != nil {
			slog.Error("executor error", "name", stage.ExecutorName, "error", err)
			break
		}
		slog.Info("executor done", "result", executorResult)
		parserResult, err := stage.Parser.Run(executorResult, stage.ParserConfig)
		if err != nil {
			slog.Error("parser error", "name", stage.ExecutorName, "error", err)
			break
		}
		slog.Info("parser done", "result", parserResult)
		parserResults = append(parserResults, StageResult{
			Name:         stage.Name,
			ParserResult: parserResult,
		})
	}
	return parserResults
}

func Cleanup() {
	for _, executor := range executorMap {
		executor.Cleanup()
	}
}
