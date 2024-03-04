package dummy

import (
	"focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/internal/stage"
	"github.com/criyle/go-judge/envexec"
)

type Dummy struct{}

func (e *Dummy) Run(stage.Cmd) (*stage.Result, error) {
	return &stage.Result{
		Status:     stage.Status(envexec.StatusInvalid),
		ExitStatus: 0,
		Error:      "I'm a dummy",
		Time:       0,
		Memory:     0,
		RunTime:    0,
		Files:      map[string]string{},
		FileIDs:    map[string]string{},
	}, nil
}

func (e *Dummy) Cleanup() error {
	return nil
}
