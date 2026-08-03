package sandbox

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/criyle/go-judge/pb"
	"github.com/joint-online-judge/JOJ3/internal/stage"
)

func TestConvertPBFileVariants(t *testing.T) {
	content, fileID, name := "content", "cached", "output"
	max := int64(123)
	tests := []struct {
		name  string
		file  stage.CmdFile
		check func(*pb.Request_File) bool
	}{
		{"content", stage.CmdFile{Content: &content}, func(f *pb.Request_File) bool { return string(f.GetMemory().GetContent()) == content }},
		{"cached", stage.CmdFile{FileID: &fileID}, func(f *pb.Request_File) bool { return f.GetCached().GetFileID() == fileID }},
		{"pipe", stage.CmdFile{Name: &name, Max: &max, Pipe: true}, func(f *pb.Request_File) bool {
			return f.GetPipe().GetName() == name && f.GetPipe().GetMax() == max && f.GetPipe().GetPipe()
		}},
		{"stream-in", stage.CmdFile{StreamIn: true}, func(f *pb.Request_File) bool { return f.GetStreamIn() != nil }},
		{"stream-out", stage.CmdFile{StreamOut: true}, func(f *pb.Request_File) bool { return f.GetStreamOut() != nil }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertPBFile(tt.file)
			if err != nil || got == nil || !tt.check(got) {
				t.Fatalf("convertPBFile() = %v, %v", got, err)
			}
		})
	}
}

func TestConvertCopyOutAndResult(t *testing.T) {
	copyOut := convertPBCopyOut([]string{"required", "optional?"})
	if len(copyOut) != 2 || copyOut[0].GetOptional() || !copyOut[1].GetOptional() || copyOut[1].GetName() != "optional" {
		t.Fatalf("convertPBCopyOut() = %+v", copyOut)
	}
	fileError := &pb.Response_FileError{}
	fileError.SetName("input")
	fileError.SetType(pb.Response_FileError_ErrorType(1))
	fileError.SetMessage("bad file")
	response := &pb.Response_Result{}
	response.SetStatus(pb.Response_Result_StatusType(stage.StatusAccepted))
	response.SetFiles(map[string][]byte{"stdout": []byte("ok")})
	response.SetFileIDs(map[string]string{"bin": "id"})
	response.SetFileError([]*pb.Response_FileError{fileError})
	got := convertPBResult([]*pb.Response_Result{response})
	if len(got) != 1 || got[0].Files["stdout"] != "ok" || got[0].FileIDs["bin"] != "id" ||
		len(got[0].FileError) != 1 || got[0].FileError[0].Message != "bad file" {
		t.Fatalf("convertPBResult() = %+v", got)
	}
}

func TestConvertPBCmdReturnsSourceReadError(t *testing.T) {
	missing := t.TempDir() + "/missing"
	_, err := convertPBCmd([]stage.Cmd{{
		Args: []string{"true"},
		CopyIn: map[string]stage.CmdFile{
			"input": {Src: &missing},
		},
	}})
	if err == nil || !strings.Contains(err.Error(), "read source file") {
		t.Fatalf("convertPBCmd() error = %v, want source read error", err)
	}
}

func TestConvertPBCmdReturnsStandardFileReadError(t *testing.T) {
	missing := t.TempDir() + "/missing"
	_, err := convertPBCmd([]stage.Cmd{{
		Args:  []string{"true"},
		Stdin: &stage.CmdFile{Src: &missing},
	}})
	if err == nil || !strings.Contains(err.Error(), "standard file") ||
		!strings.Contains(err.Error(), "read source file") {
		t.Fatalf("convertPBCmd() error = %v, want standard source read error", err)
	}
}

func TestConvertPBCmdReturnsCopyInDirectoryWalkError(t *testing.T) {
	missing := t.TempDir() + "/missing"
	_, err := convertPBCmd([]stage.Cmd{{
		Args:      []string{"true"},
		CopyInDir: missing,
	}})
	if err == nil || !strings.Contains(err.Error(), "walk") {
		t.Fatalf("convertPBCmd() error = %v, want directory walk error", err)
	}
}

func TestConvertPBCmdPreservesFilesAcrossMultipleSCMMaxFDBatches(t *testing.T) {
	// 1001 is large enough to require several descriptor batches while keeping
	// the fixture cheap to create and the protobuf request small.
	const fileCount = 1001
	dir := t.TempDir()
	for i := range fileCount {
		path := filepath.Join(dir, fmt.Sprintf("file-%03d", i))
		if err := os.WriteFile(path, []byte("content"), 0o600); err != nil {
			t.Fatal(err)
		}
	}

	cmds, err := convertPBCmd([]stage.Cmd{{
		Args:      []string{"true"},
		CopyIn:    make(map[string]stage.CmdFile),
		CopyInDir: dir,
	}})
	if err != nil {
		t.Fatal(err)
	}
	if len(cmds) != 1 || len(cmds[0].GetCopyIn()) != fileCount {
		t.Fatalf("converted %d commands with %d files, want 1 command with %d files",
			len(cmds), len(cmds[0].GetCopyIn()), fileCount)
	}
}
