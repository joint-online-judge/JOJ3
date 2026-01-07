package sandbox

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/criyle/go-judge/pb"
	"github.com/joint-online-judge/JOJ3/internal/stage"
	"google.golang.org/protobuf/types/known/emptypb"
)

// copied from https://github.com/criyle/go-judge/blob/master/cmd/go-judge-shell/grpc.go
func convertPBCmd(cmd []stage.Cmd) []*pb.Request_CmdType {
	ret := make([]*pb.Request_CmdType, 0, len(cmd))
	for _, c := range cmd {
		req := &pb.Request_CmdType{}
		req.SetArgs(c.Args)
		req.SetEnv(c.Env)
		req.SetTty(c.TTY)
		req.SetFiles(convertPBFiles([]*stage.CmdFile{c.Stdin, c.Stdout, c.Stderr}))
		req.SetCpuTimeLimit(c.CPULimit)
		req.SetClockTimeLimit(c.ClockLimit)
		req.SetMemoryLimit(c.MemoryLimit)
		req.SetStackLimit(c.StackLimit)
		req.SetProcLimit(c.ProcLimit)
		req.SetCpuRateLimit(c.CPURateLimit)
		req.SetCpuSetLimit(c.CPUSetLimit)
		req.SetDataSegmentLimit(c.DataSegmentLimit)
		req.SetAddressSpaceLimit(c.AddressSpaceLimit)
		req.SetCopyIn(convertPBCopyIn(c.CopyIn, c.CopyInDir))
		req.SetCopyOut(convertPBCopyOut(c.CopyOut))
		req.SetCopyOutCached(convertPBCopyOut(c.CopyOutCached))
		req.SetCopyOutMax(c.CopyOutMax)
		req.SetCopyOutDir(c.CopyOutDir)
		req.SetSymlinks(convertSymlink(c.CopyIn))
		ret = append(ret, req)
	}
	return ret
}

func convertPBCopyIn(
	copyIn map[string]stage.CmdFile, copyInDir string,
) map[string]*pb.Request_File {
	if copyInDir != "" {
		_ = filepath.Walk(copyInDir,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return nil
				}
				absPath, err := filepath.Abs(path)
				if err != nil {
					return nil
				}
				relPath, err := filepath.Rel(copyInDir, path)
				if err != nil {
					return nil
				}
				_, exists := copyIn[relPath]
				if !info.IsDir() && !exists {
					copyIn[relPath] = stage.CmdFile{Src: &absPath}
				}
				return nil
			})
	}
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
	rt := make([]*pb.Request_CmdCopyOutFile, 0, len(copyOut))
	for _, n := range copyOut {
		optional := false
		if strings.HasSuffix(n, "?") {
			optional = true
			n = strings.TrimSuffix(n, "?")
		}
		elem := &pb.Request_CmdCopyOutFile{}
		elem.SetName(n)
		elem.SetOptional(optional)
		rt = append(rt, elem)
	}
	return rt
}

func convertSymlink(copyIn map[string]stage.CmdFile) map[string]string {
	ret := make(map[string]string)
	for k, v := range copyIn {
		if v.Symlink == nil {
			continue
		}
		ret[k] = *v.Symlink
	}
	return ret
}

func convertPBFiles(files []*stage.CmdFile) []*pb.Request_File {
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

func convertPBFile(i stage.CmdFile) *pb.Request_File {
	req := &pb.Request_File{}
	switch {
	case i.Src != nil:
		if !filepath.IsAbs(*i.Src) {
			absPath, err := filepath.Abs(*i.Src)
			if err != nil {
				slog.Error("convert pb file get abs path", "path", *i.Src, "error", err)
				absPath = "/"
			}
			i.Src = &absPath
		}
		s, err := os.ReadFile(*i.Src)
		if err != nil {
			s = []byte{}
			slog.Error("convert pb file read file", "path", *i.Src, "error", err)
		}
		m := &pb.Request_MemoryFile{}
		m.SetContent(s)
		req.SetMemory(m)
		return req
	case i.Content != nil:
		s := strToBytes(*i.Content)
		m := &pb.Request_MemoryFile{}
		m.SetContent(s)
		req.SetMemory(m)
		return req
	case i.FileID != nil:
		c := &pb.Request_CachedFile{}
		c.SetFileID(*i.FileID)
		req.SetCached(c)
		return req
	case i.Name != nil && i.Max != nil:
		p := &pb.Request_PipeCollector{}
		p.SetName(*i.Name)
		p.SetMax(*i.Max)
		p.SetPipe(i.Pipe)
		req.SetPipe(p)
		return req
	case i.StreamIn:
		req.SetStreamIn(&emptypb.Empty{})
		return req
	case i.StreamOut:
		req.SetStreamOut(&emptypb.Empty{})
		return req
	}
	return nil
}

func convertPBResult(res []*pb.Response_Result) []stage.ExecutorResult {
	ret := make([]stage.ExecutorResult, 0, len(res))
	for _, r := range res {
		ret = append(ret, stage.ExecutorResult{
			Status:     stage.Status(r.GetStatus()),
			ExitStatus: int(r.GetExitStatus()),
			Error:      r.GetError(),
			Time:       r.GetTime(),
			Memory:     r.GetMemory(),
			RunTime:    r.GetRunTime(),
			ProcPeak:   r.GetProcPeak(),
			Files:      convertFiles(r.GetFiles()),
			Buffs:      r.GetFiles(),
			FileIDs:    r.GetFileIDs(),
			FileError:  convertPBFileError(r.GetFileError()),
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

func convertPBFileError(fe []*pb.Response_FileError) []stage.FileError {
	ret := make([]stage.FileError, 0, len(fe))
	for _, v := range fe {
		ret = append(ret, stage.FileError{
			Name:    v.GetName(),
			Type:    stage.FileErrorType(v.GetType()),
			Message: v.GetMessage(),
		})
	}
	return ret
}
