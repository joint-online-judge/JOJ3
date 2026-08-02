package sandbox

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"maps"
	"math"
	"time"

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
	callCtx, cancel := context.WithTimeout(ctx, execRPCTimeout(cmds))
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
		callCtx, cancel := context.WithTimeout(ctx, rpcTimeoutMargin)
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

func execRPCTimeout(cmds []stage.Cmd) time.Duration {
	var maxLimit uint64
	for _, cmd := range cmds {
		limit := cmd.ClockLimit
		if limit == 0 {
			// Match the local executor's default wall-clock allowance when only a
			// CPU limit is specified.
			if cmd.CPULimit > math.MaxUint64/2 {
				limit = math.MaxUint64
			} else {
				limit = cmd.CPULimit * 2
			}
		}
		if limit > maxLimit {
			maxLimit = limit
		}
	}

	margin := uint64(rpcTimeoutMargin)
	if maxLimit > uint64(math.MaxInt64)-margin {
		return time.Duration(math.MaxInt64)
	}
	// The bound above guarantees the conversion fits in time.Duration.
	return time.Duration(maxLimit + margin) // #nosec G115
}
