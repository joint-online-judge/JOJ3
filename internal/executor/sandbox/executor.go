package sandbox

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"maps"

	"github.com/criyle/go-judge/pb"
	"github.com/joint-online-judge/JOJ3/internal/stage"
	"google.golang.org/protobuf/proto"
)

func (e *Sandbox) Run(ctx context.Context, cmds []stage.Cmd) ([]stage.ExecutorResult, error) {
	var err error
	if e.execClient == nil {
		slog.Debug("create exec client", "server", e.execServer)
		e.execClient, e.conn, err = createExecClient(e.execServer, e.token)
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
	pbCmds, err := convertPBCmd(cmds)
	if err != nil {
		return nil, err
	}
	for i, pbCmd := range pbCmds {
		slog.Debug("sandbox execute", "i", i, "pbCmd size", proto.Size(pbCmd))
	}
	pbReq := &pb.Request{}
	pbReq.SetCmd(pbCmds)
	slog.Debug("sandbox execute", "pbReq size", proto.Size(pbReq))
	callCtx, cancel := context.WithTimeout(ctx, e.timeout)
	defer cancel()
	pbRet, err := e.execClient.Exec(callCtx, pbReq)
	if err != nil {
		return nil, err
	}
	if pbRet.GetError() != "" {
		return nil, fmt.Errorf("sandbox execute error: %s", pbRet.GetError())
	}
	results := convertPBResult(pbRet.GetResults())
	for _, result := range results {
		maps.Copy(e.cachedMap, result.FileIDs)
	}
	return results, nil
}

func (e *Sandbox) Cleanup(ctx context.Context) error {
	var cleanupErr error
	for k, fileID := range e.cachedMap {
		req := &pb.FileID{}
		req.SetFileID(fileID)
		callCtx, cancel := context.WithTimeout(ctx, e.timeout)
		_, err := e.execClient.FileDelete(callCtx, req)
		cancel()
		if err != nil {
			slog.Error("sandbox cleanup", "error", err)
			cleanupErr = errors.Join(cleanupErr, err)
		}
		delete(e.cachedMap, k)
	}
	if e.conn != nil {
		cleanupErr = errors.Join(cleanupErr, e.conn.Close())
		e.conn = nil
		e.execClient = nil
	}
	return cleanupErr
}
