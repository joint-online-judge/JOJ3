package stage

import (
	"github.com/criyle/go-judge/cmd/go-judge/model"
)

var executorMap = map[string]Executor{}

type Executor interface {
	Run(model.Cmd) (*model.Result, error)
}

func RegisterExecutor(name string, executor Executor) {
	executorMap[name] = executor
}
