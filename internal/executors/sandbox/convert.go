package sandbox

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/criyle/go-judge/pb"
	"github.com/joint-online-judge/JOJ3/internal/stage"
)

// copied from https://github.com/criyle/go-judge/blob/master/cmd/go-judge-shell/grpc.go
func convertPBCmd(cmd []stage.Cmd) []*pb.Request_CmdType {
	var ret []*pb.Request_CmdType
	for _, c := range cmd {
		ret = append(ret, &pb.Request_CmdType{
			Args:              c.Args,
			Env:               c.Env,
			Tty:               c.TTY,
			Files:             convertPBFiles([]*stage.CmdFile{c.Stdin, c.Stdout, c.Stderr}),
			CpuTimeLimit:      c.CPULimit,
			ClockTimeLimit:    c.ClockLimit,
			MemoryLimit:       c.MemoryLimit,
			StackLimit:        c.StackLimit,
			ProcLimit:         c.ProcLimit,
			CpuRateLimit:      c.CPURateLimit,
			CpuSetLimit:       c.CPUSetLimit,
			DataSegmentLimit:  c.DataSegmentLimit,
			AddressSpaceLimit: c.AddressSpaceLimit,
			CopyIn:            convertPBCopyIn(c.CopyIn, c.CopyInDir),
			CopyOut:           convertPBCopyOut(c.CopyOut),
			CopyOutCached:     convertPBCopyOut(c.CopyOutCached),
			CopyOutMax:        c.CopyOutMax,
			CopyOutDir:        c.CopyOutDir,
			Symlinks:          convertSymlink(c.CopyIn),
		})
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
				if !info.IsDir() {
					copyIn[path] = stage.CmdFile{Src: &absPath}
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
		rt = append(rt, &pb.Request_CmdCopyOutFile{
			Name:     n,
			Optional: optional,
		})
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
	switch {
	case i.Src != nil:
		if !filepath.IsAbs(*i.Src) {
			absPath, err := filepath.Abs(*i.Src)
			if err == nil {
				i.Src = &absPath
			}
		}
		s, err := os.ReadFile(*i.Src)
		if err != nil {
			s = []byte{}
			slog.Error("read file error", "path", *i.Src, "error", err)
		}
		return &pb.Request_File{File: &pb.Request_File_Memory{Memory: &pb.Request_MemoryFile{Content: s}}}
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

func convertPBResult(res []*pb.Response_Result) []stage.ExecutorResult {
	var ret []stage.ExecutorResult
	for _, r := range res {
		ret = append(ret, stage.ExecutorResult{
			Status:     stage.Status(r.Status),
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

func convertPBFileError(fe []*pb.Response_FileError) []stage.FileError {
	var ret []stage.FileError
	for _, v := range fe {
		ret = append(ret, stage.FileError{
			Name:    v.Name,
			Type:    stage.FileErrorType(v.Type),
			Message: v.Message,
		})
	}
	return ret
}
