package sandbox

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/criyle/go-judge/pb"
	"github.com/joint-online-judge/JOJ3/internal/stage"
	"google.golang.org/protobuf/proto"
)

type Sandbox struct {
	execServer, token string
	cachedMap         map[string]string
	execClient        pb.ExecutorClient
}

func (e *Sandbox) Run(cmds []stage.Cmd) ([]stage.ExecutorResult, error) {
	var err error
	if e.execClient == nil {
		slog.Debug("create exec client", "server", e.execServer)
		e.execClient, err = createExecClient(e.execServer, e.token)
		if err != nil {
			return nil, err
		}
	}
	// cannot use range loop since we need to change the value
	for i := 0; i < len(cmds); i += 1 {
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
	pbCmds := convertPBCmd(cmds)
	for i, pbCmd := range pbCmds {
		slog.Debug("sandbox execute", "i", i, "pbCmd size", proto.Size(pbCmd))
	}
	pbReq := &pb.Request{Cmd: pbCmds}
	slog.Info("sandbox execute", "pbReq size", proto.Size(pbReq))
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
