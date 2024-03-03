package sandbox

import (
	"strings"

	"github.com/criyle/go-judge/cmd/go-judge/model"
	"github.com/criyle/go-judge/pb"
)

// copied from https://github.com/criyle/go-judge/blob/master/cmd/go-judge-shell/grpc.go
func convertPBCmd(cmd []model.Cmd) []*pb.Request_CmdType {
	var ret []*pb.Request_CmdType
	for _, c := range cmd {
		ret = append(ret, &pb.Request_CmdType{
			Args:              c.Args,
			Env:               c.Env,
			Tty:               c.TTY,
			Files:             convertPBFiles(c.Files),
			CpuTimeLimit:      c.CPULimit,
			ClockTimeLimit:    c.ClockLimit,
			MemoryLimit:       c.MemoryLimit,
			StackLimit:        c.StackLimit,
			ProcLimit:         c.ProcLimit,
			CpuRateLimit:      c.CPURateLimit,
			CpuSetLimit:       c.CPUSetLimit,
			DataSegmentLimit:  c.DataSegmentLimit,
			AddressSpaceLimit: c.AddressSpaceLimit,
			CopyIn:            convertPBCopyIn(c.CopyIn),
			CopyOut:           convertPBCopyOut(c.CopyOut),
			CopyOutCached:     convertPBCopyOut(c.CopyOutCached),
			CopyOutMax:        c.CopyOutMax,
			CopyOutDir:        c.CopyOutDir,
			Symlinks:          convertSymlink(c.CopyIn),
		})
	}
	return ret
}

func convertPBCopyIn(copyIn map[string]model.CmdFile) map[string]*pb.Request_File {
	rt := make(map[string]*pb.Request_File, len(copyIn))
	for k, i := range copyIn {
		if i.Symlink != nil {
			continue
		}
		rt[k] = convertPBFile(i)
	}
	return rt
}

func convertPBCopyOut(copyOut []string) []*pb.Request_CmdCopyOutFile {
	rt := make([]*pb.Request_CmdCopyOutFile, 0)
	for _, n := range copyOut {
		optional := false
		if strings.HasSuffix(n, "?") {
			optional = true
			n = strings.TrimSuffix(n, "?")
		}
		rt = append(rt, &pb.Request_CmdCopyOutFile{
			Name:     n,
			Optional: optional,
		})
	}
	return rt
}

func convertSymlink(copyIn map[string]model.CmdFile) map[string]string {
	ret := make(map[string]string)
	for k, v := range copyIn {
		if v.Symlink == nil {
			continue
		}
		ret[k] = *v.Symlink
	}
	return ret
}

func convertPBFiles(files []*model.CmdFile) []*pb.Request_File {
	var ret []*pb.Request_File
	for _, f := range files {
		if f == nil {
			ret = append(ret, nil)
		} else {
			ret = append(ret, convertPBFile(*f))
		}
	}
	return ret
}

func convertPBFile(i model.CmdFile) *pb.Request_File {
	switch {
	case i.Src != nil:
		return &pb.Request_File{File: &pb.Request_File_Local{Local: &pb.Request_LocalFile{Src: *i.Src}}}
	case i.Content != nil:
		s := strToBytes(*i.Content)
		return &pb.Request_File{File: &pb.Request_File_Memory{Memory: &pb.Request_MemoryFile{Content: s}}}
	case i.FileID != nil:
		return &pb.Request_File{File: &pb.Request_File_Cached{Cached: &pb.Request_CachedFile{FileID: *i.FileID}}}
	case i.Name != nil && i.Max != nil:
		return &pb.Request_File{File: &pb.Request_File_Pipe{Pipe: &pb.Request_PipeCollector{Name: *i.Name, Max: *i.Max, Pipe: i.Pipe}}}
	case i.StreamIn:
		return &pb.Request_File{File: &pb.Request_File_StreamIn{}}
	case i.StreamOut:
		return &pb.Request_File{File: &pb.Request_File_StreamOut{}}
	}
	return nil
}

func convertPBPipeMapping(pm []model.PipeMap) []*pb.Request_PipeMap {
	var ret []*pb.Request_PipeMap
	for _, p := range pm {
		ret = append(ret, &pb.Request_PipeMap{
			In:    convertPBPipeIndex(p.In),
			Out:   convertPBPipeIndex(p.Out),
			Name:  p.Name,
			Proxy: p.Proxy,
			Max:   uint64(p.Max),
		})
	}
	return ret
}

func convertPBPipeIndex(pi model.PipeIndex) *pb.Request_PipeMap_PipeIndex {
	return &pb.Request_PipeMap_PipeIndex{Index: int32(pi.Index), Fd: int32(pi.Fd)}
}

func convertPBResult(res []*pb.Response_Result) []model.Result {
	var ret []model.Result
	for _, r := range res {
		ret = append(ret, model.Result{
			Status:     model.Status(r.Status),
			ExitStatus: int(r.ExitStatus),
			Error:      r.Error,
			Time:       r.Time,
			RunTime:    r.RunTime,
			Memory:     r.Memory,
			Files:      convertFiles(r.Files),
			Buffs:      r.Files,
			FileIDs:    r.FileIDs,
			FileError:  convertPBFileError(r.FileError),
		})
	}
	return ret
}

func convertFiles(buf map[string][]byte) map[string]string {
	ret := make(map[string]string, len(buf))
	for k, v := range buf {
		ret[k] = byteArrayToString(v)
	}
	return ret
}

func convertPBRequest(req *model.Request) *pb.StreamRequest {
	ret := &pb.StreamRequest{
		Request: &pb.StreamRequest_ExecRequest{
			ExecRequest: &pb.Request{
				RequestID:   req.RequestID,
				Cmd:         convertPBCmd(req.Cmd),
				PipeMapping: convertPBPipeMapping(req.PipeMapping),
			},
		},
	}
	return ret
}

func convertPBFileError(fe []*pb.Response_FileError) []model.FileError {
	var ret []model.FileError
	for _, v := range fe {
		ret = append(ret, model.FileError{
			Name:    v.Name,
			Type:    model.FileErrorType(v.Type),
			Message: v.Message,
		})
	}
	return ret
}
