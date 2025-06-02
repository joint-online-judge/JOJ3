package stage

import (
	"encoding/json"
	"fmt"
	"strconv"
)

var executorMap = map[string]Executor{}

type Executor interface {
	Run([]Cmd) ([]ExecutorResult, error)
	Cleanup() error
}

func RegisterExecutor(name string, executor Executor) {
	executorMap[name] = executor
}

// ExecutorResult defines single command result
type ExecutorResult struct {
	Status     Status            `json:"status"`
	ExitStatus int               `json:"exitStatus"`
	Error      string            `json:"error,omitempty"`
	Time       uint64            `json:"time"`    // ns (cgroup recorded time)
	Memory     uint64            `json:"memory"`  // byte
	RunTime    uint64            `json:"runTime"` // ns (wall clock time)
	ProcPeak   uint64            `json:"procPeak"`
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
	summary.Status = StatusAccepted
	for _, result := range results {
		if result.Status != StatusAccepted &&
			summary.Status == StatusAccepted {
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
