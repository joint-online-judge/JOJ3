package stage

import (
	"github.com/criyle/go-judge/cmd/go-judge/model"
)

type Stage struct {
	name         string
	executor     Executor
	executorCmd  model.Cmd
	parser       Parser
	parserConfig string
}

type StageResult struct {
	Name string
	ParserResult
}

func ParseStages() []Stage {
	stages := []Stage{}
	config := [][]string{
		{"dummy stage 0", "dummy", "dummy"},
		{"dummy stage 1", "dummy", "dummy"},
	}
	for _, v := range config {
		stages = append(stages, Stage{
			name:         v[0],
			executor:     executorMap[v[1]],
			executorCmd:  model.Cmd{},
			parser:       parserMap[v[2]],
			parserConfig: "",
		})
	}
	return stages
}

func Run(stages []Stage) []StageResult {
	var parserResults []StageResult
	for _, stage := range stages {
		executorResult := stage.executor.Run(stage.executorCmd)
		parserResult := stage.parser.Run(executorResult, stage.parserConfig)
		parserResults = append(parserResults, StageResult{
			Name:         stage.name,
			ParserResult: parserResult,
		})
	}
	return parserResults
}
