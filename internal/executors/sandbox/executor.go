package sandbox

import (
	"context"
	"fmt"
	"log/slog"

	"focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/internal/stage"
	"github.com/criyle/go-judge/pb"
)

type Sandbox struct {
	execClient pb.ExecutorClient
}

func (e *Sandbox) Run(cmd stage.Cmd) (*stage.Result, error) {
	slog.Info("sandbox run", "cmd", cmd)
	req := &pb.Request{Cmd: convertPBCmd([]stage.Cmd{cmd})}
	ret, err := e.execClient.Exec(context.TODO(), req)
	if err != nil {
		return nil, err
	}
	if ret.Error != "" {
		return nil, fmt.Errorf("compile error: %s", ret.Error)
	}
	slog.Info("sandbox run", "ret", ret)
	return &convertPBResult(ret.Results)[0], nil
}
