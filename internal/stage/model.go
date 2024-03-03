package stage

import (
	"github.com/criyle/go-judge/cmd/go-judge/model"
)

type Stage struct {
	Name         string
	ExecutorName string
	Executor     Executor
	ExecutorCmd  model.Cmd
	ParserName   string
	Parser       Parser
	ParserConfig any
}

type ParserResult struct {
	Score   int
	Comment string
}

type StageResult struct {
	Name string
	*ParserResult
}

type StagesConfig struct {
	Stages []struct {
		Name     string
		Executor struct {
			Name string
			With model.Cmd
		}
		Parser struct {
			Name string
			With interface{}
		}
	}
}
