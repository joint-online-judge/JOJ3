package dummy

import (
	"github.com/criyle/go-judge/cmd/go-judge/model"
	"github.com/criyle/go-judge/envexec"
)

type Dummy struct{}

func (e *Dummy) Run(model.Cmd) model.Result {
	return model.Result{
		Status:     model.Status(envexec.StatusInvalid),
		ExitStatus: 0,
		Error:      "I'm a dummy",
		Time:       0,
		Memory:     0,
		RunTime:    0,
		Files:      map[string]string{},
		FileIDs:    map[string]string{},
	}
}
