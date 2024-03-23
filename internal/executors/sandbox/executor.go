package sandbox

import (
	"context"
	"fmt"
	"log/slog"

	"focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/internal/stage"
	"github.com/criyle/go-judge/pb"
)

type Sandbox struct {
	execServer, token string
	cachedMap         map[string]string
	execClient        pb.ExecutorClient
}

func (e *Sandbox) Run(cmds []stage.Cmd) ([]stage.ExecutorResult, error) {
	var err error
	if e.execClient == nil {
		e.execClient, err = createExecClient(e.execServer, e.token)
		if err != nil {
			return nil, err
		}
	}
	// cannot use range loop since we need to change the value
	for i := 0; i < len(cmds); i++ {
		cmd := &cmds[i]
		if cmd.CopyIn == nil {
			cmd.CopyIn = make(map[string]stage.CmdFile)
		}
		for k, v := range cmd.CopyInCached {
			if fileID, ok := e.cachedMap[v]; ok {
				cmd.CopyIn[k] = stage.CmdFile{FileID: &fileID}
			}
		}
	}
	pbReq := &pb.Request{Cmd: convertPBCmd(cmds)}
	pbRet, err := e.execClient.Exec(context.TODO(), pbReq)
	if err != nil {
		return nil, err
	}
	if pbRet.Error != "" {
		return nil, fmt.Errorf("sandbox execute error: %s", pbRet.Error)
	}
	results := convertPBResult(pbRet.Results)
	for _, result := range results {
		for fileName, fileID := range result.FileIDs {
			e.cachedMap[fileName] = fileID
		}
	}
	return results, nil
}

func (e *Sandbox) Cleanup() error {
	for k, fileID := range e.cachedMap {
		_, err := e.execClient.FileDelete(context.TODO(), &pb.FileID{
			FileID: fileID,
		})
		if err != nil {
			slog.Error("sandbox cleanup", "error", err)
		}
		delete(e.cachedMap, k)
	}
	return nil
}
