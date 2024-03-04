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
	cachedMap  map[string]string
}

func (e *Sandbox) Run(cmd stage.Cmd) (*stage.ExecutorResult, error) {
	slog.Info("sandbox run", "cmd", cmd)
	if cmd.CopyIn == nil {
		cmd.CopyIn = make(map[string]stage.CmdFile)
	}
	for k, v := range cmd.CopyInCached {
		if fileID, ok := e.cachedMap[v]; ok {
			cmd.CopyIn[k] = stage.CmdFile{FileID: &fileID}
		}
	}
	req := &pb.Request{Cmd: convertPBCmd([]stage.Cmd{cmd})}
	ret, err := e.execClient.Exec(context.TODO(), req)
	if err != nil {
		return nil, err
	}
	if ret.Error != "" {
		return nil, fmt.Errorf("compile error: %s", ret.Error)
	}
	slog.Info("sandbox run", "ret", ret)
	res := &convertPBResult(ret.Results)[0]
	for fileName, fileID := range res.FileIDs {
		e.cachedMap[fileName] = fileID
	}
	return res, nil
}

func (e *Sandbox) Cleanup() error {
	slog.Info("sandbox cleanup")
	for _, fileID := range e.cachedMap {
		_, err := e.execClient.FileDelete(context.TODO(), &pb.FileID{
			FileID: fileID,
		})
		if err != nil {
			slog.Error("sandbox cleanup", "error", err)
		}
	}
	return nil
}
