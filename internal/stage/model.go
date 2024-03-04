package stage

import (
	"fmt"
	"strconv"
	"time"

	"github.com/criyle/go-judge/envexec"
)

// copied from https://github.com/criyle/go-judge/blob/master/cmd/go-judge/model/model.go
// FileError defines the location, file name and the detailed message for a failed file operation
type FileError = envexec.FileError

// FileErrorType defines the location that file operation fails
type FileErrorType = envexec.FileErrorType

// CmdFile defines file from multiple source including local / memory / cached or pipe collector
type CmdFile struct {
	Src       *string `json:"src"`
	Content   *string `json:"content"`
	FileID    *string `json:"fileId"`
	Name      *string `json:"name"`
	Max       *int64  `json:"max"`
	Symlink   *string `json:"symlink"`
	StreamIn  bool    `json:"streamIn"`
	StreamOut bool    `json:"streamOut"`
	Pipe      bool    `json:"pipe"`
}

// Cmd defines command and limits to start a program using in envexec
type Cmd struct {
	Args  []string   `json:"args"`
	Env   []string   `json:"env,omitempty"`
	Files []*CmdFile `json:"files,omitempty"`

	CPULimit     uint64 `json:"cpuLimit"`
	RealCPULimit uint64 `json:"realCpuLimit"`
	ClockLimit   uint64 `json:"clockLimit"`
	MemoryLimit  uint64 `json:"memoryLimit"`
	StackLimit   uint64 `json:"stackLimit"`
	ProcLimit    uint64 `json:"procLimit"`
	CPURateLimit uint64 `json:"cpuRateLimit"`
	CPUSetLimit  string `json:"cpuSetLimit"`

	CopyIn map[string]CmdFile `json:"copyIn"`

	CopyOut       []string `json:"copyOut"`
	CopyOutCached []string `json:"copyOutCached"`
	CopyOutMax    uint64   `json:"copyOutMax"`
	CopyOutDir    string   `json:"copyOutDir"`

	TTY               bool `json:"tty,omitempty"`
	StrictMemoryLimit bool `json:"strictMemoryLimit"`
	DataSegmentLimit  bool `json:"dataSegmentLimit"`
	AddressSpaceLimit bool `json:"addressSpaceLimit"`
}

// PipeIndex defines indexing for a pipe fd
type PipeIndex struct {
	Index int `json:"index"`
	Fd    int `json:"fd"`
}

// PipeMap defines in / out pipe for multiple program
type PipeMap struct {
	In    PipeIndex `json:"in"`
	Out   PipeIndex `json:"out"`
	Name  string    `json:"name"`
	Max   int64     `json:"max"`
	Proxy bool      `json:"proxy"`
}

// Request defines single worker request
type Request struct {
	RequestID   string    `json:"requestId"`
	Cmd         []Cmd     `json:"cmd"`
	PipeMapping []PipeMap `json:"pipeMapping"`
}

// Status offers JSON marshal for envexec.Status
type Status envexec.Status

// String converts status to string
func (s Status) String() string {
	return envexec.Status(s).String()
}

// MarshalJSON convert status into string
func (s Status) MarshalJSON() ([]byte, error) {
	return []byte("\"" + envexec.Status(s).String() + "\""), nil
}

// UnmarshalJSON convert string into status
func (s *Status) UnmarshalJSON(b []byte) error {
	str := string(b)
	v, err := envexec.StringToStatus(str)
	if err != nil {
		return err
	}
	*s = Status(v)
	return nil
}

// Result defines single command result
type Result struct {
	Status     Status            `json:"status"`
	ExitStatus int               `json:"exitStatus"`
	Error      string            `json:"error,omitempty"`
	Time       uint64            `json:"time"`
	Memory     uint64            `json:"memory"`
	RunTime    uint64            `json:"runTime"`
	Files      map[string]string `json:"files,omitempty"`
	FileIDs    map[string]string `json:"fileIds,omitempty"`
	FileError  []FileError       `json:"fileError,omitempty"`

	files []string
	Buffs map[string][]byte `json:"-"`
}

func (r Result) String() string {
	type Result struct {
		Status     Status
		ExitStatus int
		Error      string
		Time       time.Duration
		RunTime    time.Duration
		Memory     envexec.Size
		Files      map[string]string
		FileIDs    map[string]string
		FileError  []FileError
	}
	d := Result{
		Status:     r.Status,
		ExitStatus: r.ExitStatus,
		Error:      r.Error,
		Time:       time.Duration(r.Time),
		RunTime:    time.Duration(r.RunTime),
		Memory:     envexec.Size(r.Memory),
		Files:      make(map[string]string),
		FileIDs:    r.FileIDs,
		FileError:  r.FileError,
	}
	for k, v := range r.Files {
		d.Files[k] = "len:" + strconv.Itoa(len(v))
	}
	return fmt.Sprintf("%+v", d)
}

type Stage struct {
	Name         string
	ExecutorName string
	Executor     Executor
	ExecutorCmd  Cmd
	ParserName   string
	Parser       Parser
	ParserConfig any
}

type ParserResult struct {
	Score   int
	Comment string
}

type StageResult struct {
	Name string
	*ParserResult
}

type StagesConfig struct {
	Stages []struct {
		Name     string
		Executor struct {
			Name string
			With Cmd
		}
		Parser struct {
			Name string
			With interface{}
		}
	}
}
