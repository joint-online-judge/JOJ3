package stage

import (
	"log/slog"
	"os"

	"github.com/pelletier/go-toml/v2"
)

func ParseStages(tomlConfig []byte) []Stage {
	var stagesConfig StagesConfig
	err := toml.Unmarshal(tomlConfig, &stagesConfig)
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
		slog.Info("executor run start", "cmd", stage.ExecutorCmd)
		executorResult, err := stage.Executor.Run(stage.ExecutorCmd)
		if err != nil {
			slog.Error("executor run error", "name", stage.ExecutorName, "error", err)
			break
		}
		slog.Info("executor run done", "result", executorResult)
		slog.Info("parser run start", "config", stage.ParserConfig)
		parserResult, err := stage.Parser.Run(executorResult, stage.ParserConfig)
		if err != nil {
			slog.Error("parser run error", "name", stage.ExecutorName, "error", err)
			break
		}
		slog.Info("parser run done", "result", parserResult)
		parserResults = append(parserResults, StageResult{
			Name:         stage.Name,
			ParserResult: parserResult,
		})
	}
	return parserResults
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
