// Package resultdetail provides detailed execution result output.
package resultdetail

import "github.com/joint-online-judge/JOJ3/internal/stage"

var name = "result-detail"

type Conf struct {
	Score              int
	ShowExecutorStatus bool `default:"true"`
	ShowExitStatus     bool `default:"false"`
	ShowError          bool `default:"false"`
	ShowTime           bool `default:"true"`
	ShowMemory         bool `default:"true"`
	ShowRunTime        bool `default:"false"`
	ShowProcPeak       bool `default:"false"`
	ShowFiles          []string
	FilesInCodeBlock   bool `default:"true"`
	MaxFileLength      int  `default:"2048"`
}

type ResultDetail struct{}

func init() {
	stage.RegisterParser(name, &ResultDetail{})
}
