package sandbox

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"maps"
	"os"
	"path/filepath"

	"github.com/criyle/go-judge/pb"
	"github.com/joint-online-judge/JOJ3/internal/stage"
	"google.golang.org/protobuf/proto"
)

const tarThreshold = 128 * 1024

func (e *Sandbox) Run(cmds []stage.Cmd) ([]stage.ExecutorResult, error) {
	var err error
	if e.execClient == nil {
		slog.Debug("create exec client", "server", e.execServer)
		e.execClient, err = createExecClient(e.execServer, e.token)
		if err != nil {
			return nil, err
		}
	}
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
	if needTar, tarData := prepareTar(cmds); needTar {
		return e.runWithTar(cmds, tarData)
	}
	return e.runUnary(cmds)
}

func prepareTar(cmds []stage.Cmd) (bool, []byte) {
	if len(cmds) == 0 {
		return false, nil
	}
	for i := range cmds {
		if cmds[i].CopyInDir != "" &&
			estimateCopyInSize(&cmds[i]) >= tarThreshold {
			tarData, keysInTar := createCopyInTar(&cmds[i])
			if tarData == nil {
				return false, nil
			}
			slog.Debug("prepareTar", "tarSize", len(tarData), "keysInTar", len(keysInTar))
			for j := range cmds {
				cmds[j].CopyInDir = ""
				for _, k := range keysInTar {
					delete(cmds[j].CopyIn, k)
				}
			}
			return true, tarData
		}
	}
	return false, nil
}

func (e *Sandbox) runUnary(cmds []stage.Cmd) ([]stage.ExecutorResult, error) {
	pbCmds := convertPBCmd(cmds)
	for i, pbCmd := range pbCmds {
		slog.Debug("sandbox execute", "i", i, "pbCmd size", proto.Size(pbCmd))
	}
	pbReq := &pb.Request{}
	pbReq.SetCmd(pbCmds)
	slog.Debug("sandbox execute", "pbReq size", proto.Size(pbReq))
	pbRet, err := e.execClient.Exec(context.TODO(), pbReq)
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

func (e *Sandbox) runWithTar(cmds []stage.Cmd, tarData []byte) ([]stage.ExecutorResult, error) {
	fc := &pb.FileContent{}
	fc.SetContent(tarData)
	fileIDResp, err := e.execClient.FileAdd(context.TODO(), fc)
	if err != nil {
		return nil, fmt.Errorf("file add tar: %w", err)
	}
	fid := fileIDResp.GetFileID()
	slog.Debug("sandbox tar uploaded", "fileID", fid, "tarSize", len(tarData))

	tarFileName := "/w/__joj3_copyin.tar"
	script := fmt.Sprintf(
		"/bin/tar xf %s -C / --no-same-owner && rm %s && exec \"$@\"",
		tarFileName, tarFileName,
	)

	cmds[0].CopyIn[tarFileName] = stage.CmdFile{FileID: &fid}
	cmds[0].Args = append([]string{
		"/bin/sh", "-c", script, "--",
	}, cmds[0].Args...)

	slog.Debug("sandbox tar exec", "cmd", cmds[0].Args[:3])

	results, err := e.runUnary(cmds)
	if err != nil {
		return nil, err
	}

	deleteReq := &pb.FileID{}
	deleteReq.SetFileID(fid)
	if _, err := e.execClient.FileDelete(context.TODO(), deleteReq); err != nil {
		slog.Warn("sandbox tar file delete", "fileID", fid, "error", err)
	}

	return results, nil
}

func estimateCopyInSize(cmd *stage.Cmd) int {
	total := 0
	if cmd.CopyInDir != "" {
		_ = filepath.Walk(cmd.CopyInDir,
			func(path string, info os.FileInfo, err error) error {
				if err != nil || info.IsDir() {
					return nil
				}
				relPath, err := filepath.Rel(cmd.CopyInDir, path)
				if err != nil {
					return nil
				}
				if _, exists := cmd.CopyIn[relPath]; !exists {
					total += int(info.Size())
				}
				return nil
			})
	}
	for _, f := range cmd.CopyIn {
		if f.Symlink != nil {
			continue
		}
		if f.Src != nil {
			if fi, err := os.Stat(*f.Src); err == nil {
				total += int(fi.Size())
			}
		} else if f.Content != nil {
			total += len(*f.Content)
		}
	}
	return total
}

func createCopyInTar(cmd *stage.Cmd) ([]byte, []string) {
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	tarKeys := make([]string, 0, len(cmd.CopyIn))

	if cmd.CopyInDir != "" {
		err := filepath.Walk(cmd.CopyInDir,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				relPath, err := filepath.Rel(cmd.CopyInDir, path)
				if err != nil {
					return err
				}
				if relPath == "." {
					return nil
				}
				if _, exists := cmd.CopyIn[relPath]; exists {
					if info.IsDir() {
						return filepath.SkipDir
					}
					return nil
				}
				if info.Mode()&os.ModeSymlink != 0 {
					link, err := os.Readlink(path)
					if err != nil {
						return err
					}
					hdr, err := tar.FileInfoHeader(info, link)
					if err != nil {
						return err
					}
					hdr.Name = relPath
					if !filepath.IsAbs(hdr.Name) {
						hdr.Name = "w/" + hdr.Name
					}
					return tw.WriteHeader(hdr)
				}
				hdr, err := tar.FileInfoHeader(info, "")
				if err != nil {
					return err
				}
				hdr.Name = relPath
				if !filepath.IsAbs(hdr.Name) {
					hdr.Name = "w/" + hdr.Name
				}
				if info.IsDir() {
					hdr.Name += "/"
				}
				if err := tw.WriteHeader(hdr); err != nil {
					return err
				}
				if info.IsDir() {
					return nil
				}
				f, err := os.Open(path)
				if err != nil {
					return err
				}
				_, err = io.Copy(tw, f)
				f.Close()
				return err
			})
		if err != nil {
			slog.Error("create copyIn tar walk", "error", err)
			return nil, nil
		}
	}

	for k, f := range cmd.CopyIn {
		if f.FileID != nil || f.Symlink != nil || f.StreamIn || f.StreamOut || f.Pipe {
			continue
		}
		if f.Content != nil {
			name := k
			if !filepath.IsAbs(name) {
				name = "w/" + name
			}
			hdr := &tar.Header{
				Name: name,
				Mode: 0o644,
				Size: int64(len(*f.Content)),
			}
			if err := tw.WriteHeader(hdr); err != nil {
				slog.Error("create copyIn tar write header", "key", k, "error", err)
				continue
			}
			if _, err := tw.Write([]byte(*f.Content)); err != nil {
				slog.Error("create copyIn tar write content", "key", k, "error", err)
				continue
			}
			tarKeys = append(tarKeys, k)
		} else if f.Src != nil {
			fi, err := os.Stat(*f.Src)
			if err != nil {
				slog.Error("create copyIn tar stat", "key", k, "src", *f.Src, "error", err)
				continue
			}
			if fi.IsDir() {
				continue
			}
			hdr, err := tar.FileInfoHeader(fi, "")
			if err != nil {
				slog.Error("create copyIn tar file info header", "key", k, "error", err)
				continue
			}
			hdr.Name = k
			if !filepath.IsAbs(hdr.Name) {
				hdr.Name = "w/" + hdr.Name
			}
			if err := tw.WriteHeader(hdr); err != nil {
				slog.Error("create copyIn tar write header", "key", k, "error", err)
				continue
			}
			srcFile, err := os.Open(*f.Src)
			if err != nil {
				slog.Error("create copyIn tar open src", "key", k, "src", *f.Src, "error", err)
				continue
			}
			_, err = io.Copy(tw, srcFile)
			srcFile.Close()
			if err != nil {
				slog.Error("create copyIn tar copy", "key", k, "error", err)
				continue
			}
			tarKeys = append(tarKeys, k)
		}
	}

	if err := tw.Close(); err != nil {
		slog.Error("create copyIn tar close", "error", err)
		return nil, nil
	}
	return buf.Bytes(), tarKeys
}

func (e *Sandbox) Cleanup() error {
	for k, fileID := range e.cachedMap {
		req := &pb.FileID{}
		req.SetFileID(fileID)
		_, err := e.execClient.FileDelete(context.TODO(), req)
		if err != nil {
			slog.Error("sandbox cleanup", "error", err)
		}
		delete(e.cachedMap, k)
	}
	return nil
}
