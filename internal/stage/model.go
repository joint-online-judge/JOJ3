package stage

import (
	"encoding/json"
	"fmt"
	"strconv"

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
	Args   []string `json:"args"`
	Env    []string `json:"env,omitempty"`
	Stdin  *CmdFile `json:"stdin,omitempty"`
	Stdout *CmdFile `json:"stdout,omitempty"`
	Stderr *CmdFile `json:"stderr,omitempty"`

	CPULimit     uint64 `json:"cpuLimit"`     // ns
	RealCPULimit uint64 `json:"realCpuLimit"` // deprecated: use clock limit instead (still working)
	ClockLimit   uint64 `json:"clockLimit"`   // ns
	MemoryLimit  uint64 `json:"memoryLimit"`  // byte
	StackLimit   uint64 `json:"stackLimit"`   // byte
	ProcLimit    uint64 `json:"procLimit"`
	CPURateLimit uint64 `json:"cpuRateLimit"` // limit cpu usage (1000 equals 1 cpu)
	CPUSetLimit  string `json:"cpuSetLimit"`  // set the cpuSet for cgroup

	CopyIn       map[string]CmdFile `json:"copyIn"`
	CopyInCached map[string]string  `json:"copyInCached"`
	CopyInDir    string             `json:"copyInDir"`

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

// ExecutorResult defines single command result
type ExecutorResult struct {
	Status     Status            `json:"status"`
	ExitStatus int               `json:"exitStatus"`
	Error      string            `json:"error,omitempty"`
	Time       uint64            `json:"time"`    // ns (cgroup recorded time)
	Memory     uint64            `json:"memory"`  // byte
	RunTime    uint64            `json:"runTime"` // ns (wall clock time)
	Files      map[string]string `json:"files,omitempty"`
	FileIDs    map[string]string `json:"fileIds,omitempty"`
	FileError  []FileError       `json:"fileError,omitempty"`

	Buffs map[string][]byte `json:"-"`
}

type ExecutorResultSummary struct {
	Status     Status            `json:"status"`
	ExitStatus int               `json:"exitStatus"`
	Error      string            `json:"error,omitempty"`
	Time       uint64            `json:"time"`    // ns (cgroup recorded time)
	Memory     uint64            `json:"memory"`  // byte
	RunTime    uint64            `json:"runTime"` // ns (wall clock time)
	Files      map[string]string `json:"files,omitempty"`
	FileIDs    map[string]string `json:"fileIds,omitempty"`
	FileError  []FileError       `json:"fileError,omitempty"`
}

func (r ExecutorResult) String() string {
	d := ExecutorResultSummary{
		Status:     r.Status,
		ExitStatus: r.ExitStatus,
		Error:      r.Error,
		Time:       r.Time,
		Memory:     r.Memory,
		RunTime:    r.RunTime,
		Files:      make(map[string]string),
		FileIDs:    r.FileIDs,
		FileError:  r.FileError,
	}
	for k, v := range r.Files {
		d.Files[k] = "len:" + strconv.Itoa(len(v))
	}
	return fmt.Sprintf("%+v", d)
}

func (r ExecutorResult) MarshalJSON() ([]byte, error) {
	d := ExecutorResultSummary{
		Status:     r.Status,
		ExitStatus: r.ExitStatus,
		Error:      r.Error,
		Time:       r.Time,
		Memory:     r.Memory,
		RunTime:    r.RunTime,
		Files:      make(map[string]string),
		FileIDs:    r.FileIDs,
		FileError:  r.FileError,
	}
	for k, v := range r.Files {
		d.Files[k] = "len:" + strconv.Itoa(len(v))
	}
	return json.Marshal(d)
}

func SummarizeExecutorResults(results []ExecutorResult) ExecutorResultSummary {
	var summary ExecutorResultSummary
	summary.Status = Status(envexec.StatusAccepted)
	for _, result := range results {
		if result.Status != Status(envexec.StatusAccepted) &&
			summary.Status == Status(envexec.StatusAccepted) {
			summary.Status = result.Status
		}
		if result.ExitStatus != 0 && summary.ExitStatus == 0 {
			summary.ExitStatus = result.ExitStatus
		}
		if result.Error != "" && summary.Error == "" {
			summary.Error = result.Error
		}
		summary.Time += result.Time
		summary.Memory += result.Memory
		summary.RunTime += result.RunTime
	}
	return summary
}

type StageExecutor struct {
	Name string
	Cmds []Cmd
}
type StageParser struct {
	Name string
	Conf any
}

type Stage struct {
	Name     string
	Executor StageExecutor
	Parsers  []StageParser
}

type ParserResult struct {
	Score   int    `json:"score"`
	Comment string `json:"comment"`
}

type NonNullSlice[T any] []T

func (s NonNullSlice[T]) MarshalJSON() ([]byte, error) {
	if len(s) == 0 {
		return []byte("[]"), nil
	}
	return json.Marshal([]T(s))
}

type StageResult struct {
	Name      string                     `json:"name"`
	Results   NonNullSlice[ParserResult] `json:"results"`
	ForceQuit bool                       `json:"force_quit"` // underscore as it will dump to file
}

type CaseDetail struct {
	Index          int            `json:"index"`
	ExecutorResult ExecutorResult `json:"executorResult"`
	ParserScores   map[string]int `json:"parserScores"`
}

type StageDetail struct {
	Name        string       `json:"name"`
	CaseDetails []CaseDetail `json:"caseDetails"`
}
