package stage

import (
	"github.com/criyle/go-judge/cmd/go-judge/model"
	"github.com/pelletier/go-toml/v2"
)

func ParseStages(tomlConfig string) []Stage {
	var stagesConfig StagesConfig
	err := toml.Unmarshal([]byte(tomlConfig), &stagesConfig)
	if err != nil {
		panic(err)
	}
	stages := []Stage{}
	for _, stage := range stagesConfig.Stages {
		stages = append(stages, Stage{
			Name:         stage.Name,
			Executor:     executorMap[stage.Executor.Name],
			ExecutorCmd:  model.Cmd{},
			Parser:       parserMap[stage.Parser.Name],
			ParserConfig: stage.Parser.With,
		})
	}
	return stages
}

func Run(stages []Stage) []StageResult {
	var parserResults []StageResult
	for _, stage := range stages {
		executorResult := stage.Executor.Run(stage.ExecutorCmd)
		parserResult := stage.Parser.Run(executorResult, stage.ParserConfig)
		parserResults = append(parserResults, StageResult{
			Name:         stage.Name,
			ParserResult: parserResult,
		})
	}
	return parserResults
}
