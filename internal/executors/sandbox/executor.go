package sandbox

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/criyle/go-judge/cmd/go-judge/model"
	"github.com/criyle/go-judge/pb"
)

type Sandbox struct {
	execClient pb.ExecutorClient
}

func (e *Sandbox) Run(cmd model.Cmd) (*model.Result, error) {
	slog.Info("sandbox run", "cmd", cmd)
	req := &pb.Request{Cmd: convertPBCmd([]model.Cmd{cmd})}
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
