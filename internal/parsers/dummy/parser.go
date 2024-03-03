package dummy

import (
	"focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/internal/stage"
	"github.com/criyle/go-judge/cmd/go-judge/model"
	"github.com/mitchellh/mapstructure"
)

type Config struct {
	Score   int
	Comment string
}

type Dummy struct{}

func (e *Dummy) Run(result model.Result, configAny any) stage.ParserResult {
	var config Config
	err := mapstructure.Decode(configAny, &config)
	if err != nil {
		panic(err)
	}
	return stage.ParserResult{
		Score:   config.Score,
		Comment: config.Comment,
	}
}
