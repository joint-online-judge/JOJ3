package dummy

import (
	"focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/internal/stage"
	"github.com/criyle/go-judge/cmd/go-judge/model"
)

type Dummy struct{}

func (e *Dummy) Run(result model.Result, config string) stage.ParserResult {
	return stage.ParserResult{
		Score:   0,
		Comment: "I'm a dummy",
	}
}
