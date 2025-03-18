package dummy

import "github.com/joint-online-judge/JOJ3/internal/stage"

func (e *Dummy) Run(cmds []stage.Cmd) ([]stage.ExecutorResult, error) {
	var res []stage.ExecutorResult
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

func (e *Dummy) Cleanup() error {
	return nil
}
