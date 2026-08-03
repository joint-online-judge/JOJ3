package sandbox

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/criyle/go-judge/pb"
	"github.com/joint-online-judge/JOJ3/internal/stage"
	"google.golang.org/protobuf/types/known/emptypb"
)

// copied from https://github.com/criyle/go-judge/blob/master/cmd/go-judge-shell/grpc.go
func convertPBCmd(cmd []stage.Cmd) ([]*pb.Request_CmdType, error) {
	ret := make([]*pb.Request_CmdType, 0, len(cmd))
	for index, c := range cmd {
		files, err := convertPBFiles([]*stage.CmdFile{c.Stdin, c.Stdout, c.Stderr})
		if err != nil {
			return nil, fmt.Errorf("command %d standard file: %w", index, err)
		}
		copyIn, err := convertPBCopyIn(c.CopyIn, c.CopyInDir)
		if err != nil {
			return nil, fmt.Errorf("command %d copy-in: %w", index, err)
		}
		req := &pb.Request_CmdType{}
		req.SetArgs(c.Args)
		req.SetEnv(c.Env)
		req.SetTty(c.TTY)
		req.SetFiles(files)
		req.SetCpuTimeLimit(c.CPULimit)
		req.SetClockTimeLimit(c.ClockLimit)
		req.SetMemoryLimit(c.MemoryLimit)
		req.SetStackLimit(c.StackLimit)
		req.SetProcLimit(c.ProcLimit)
		req.SetCpuRateLimit(c.CPURateLimit)
		req.SetCpuSetLimit(c.CPUSetLimit)
		req.SetDataSegmentLimit(c.DataSegmentLimit)
		req.SetAddressSpaceLimit(c.AddressSpaceLimit)
		req.SetCopyIn(copyIn)
		req.SetCopyOut(convertPBCopyOut(c.CopyOut))
		req.SetCopyOutCached(convertPBCopyOut(c.CopyOutCached))
		req.SetCopyOutMax(c.CopyOutMax)
		req.SetCopyOutDir(c.CopyOutDir)
		req.SetSymlinks(convertSymlink(c.CopyIn))
		ret = append(ret, req)
	}
	return ret, nil
}

func convertPBCopyIn(
	copyIn map[string]stage.CmdFile, copyInDir string,
) (map[string]*pb.Request_File, error) {
	if copyInDir != "" {
		err := filepath.Walk(copyInDir,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				absPath, err := filepath.Abs(path)
				if err != nil {
					return err
				}
				relPath, err := filepath.Rel(copyInDir, path)
				if err != nil {
					return err
				}
				_, exists := copyIn[relPath]
				if !info.IsDir() && !exists {
					copyIn[relPath] = stage.CmdFile{Src: &absPath}
				}
				return nil
			})
		if err != nil {
			return nil, fmt.Errorf("walk %q: %w", copyInDir, err)
		}
	}
	rt := make(map[string]*pb.Request_File, len(copyIn))
	for k, i := range copyIn {
		if i.Symlink != nil {
			continue
		}
		file, err := convertPBFile(i)
		if err != nil {
			return nil, fmt.Errorf("file %q: %w", k, err)
		}
		rt[k] = file
	}
	return rt, nil
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

func convertPBFiles(files []*stage.CmdFile) ([]*pb.Request_File, error) {
	var ret []*pb.Request_File
	for _, f := range files {
		if f == nil {
			ret = append(ret, nil)
		} else {
			file, err := convertPBFile(*f)
			if err != nil {
				return nil, err
			}
			ret = append(ret, file)
		}
	}
	return ret, nil
}

func convertPBFile(i stage.CmdFile) (*pb.Request_File, error) {
	req := &pb.Request_File{}
	switch {
	case i.Src != nil:
		if !filepath.IsAbs(*i.Src) {
			absPath, err := filepath.Abs(*i.Src)
			if err != nil {
				return nil, fmt.Errorf("resolve source path %q: %w", *i.Src, err)
			}
			i.Src = &absPath
		}
		s, err := os.ReadFile(*i.Src)
		if err != nil {
			return nil, fmt.Errorf("read source file %q: %w", *i.Src, err)
		}
		m := &pb.Request_MemoryFile{}
		m.SetContent(s)
		req.SetMemory(m)
		return req, nil
	case i.Content != nil:
		s := strToBytes(*i.Content)
		m := &pb.Request_MemoryFile{}
		m.SetContent(s)
		req.SetMemory(m)
		return req, nil
	case i.FileID != nil:
		c := &pb.Request_CachedFile{}
		c.SetFileID(*i.FileID)
		req.SetCached(c)
		return req, nil
	case i.Name != nil && i.Max != nil:
		p := &pb.Request_PipeCollector{}
		p.SetName(*i.Name)
		p.SetMax(*i.Max)
		p.SetPipe(i.Pipe)
		req.SetPipe(p)
		return req, nil
	case i.StreamIn:
		req.SetStreamIn(&emptypb.Empty{})
		return req, nil
	case i.StreamOut:
		req.SetStreamOut(&emptypb.Empty{})
		return req, nil
	}
	return nil, nil
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
