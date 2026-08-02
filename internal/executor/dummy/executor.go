package dummy

import (
	"context"

	"github.com/joint-online-judge/JOJ3/internal/stage"
)

func (e *Dummy) Run(_ context.Context, cmds []stage.Cmd) ([]stage.ExecutorResult, error) {
	res := make([]stage.ExecutorResult, 0, len(cmds))
	for range cmds {
		res = append(res, stage.ExecutorResult{
			Status:     stage.StatusAccepted,
			ExitStatus: 0,
			Error:      "",
			Time:       0,
			Memory:     0,
			RunTime:    0,
			Files:      map[string]string{},
			FileIDs:    map[string]string{},
		})
	}
	return res, nil
}

func (e *Dummy) Cleanup(_ context.Context) error {
	return nil
}
