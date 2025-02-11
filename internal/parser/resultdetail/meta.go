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
	ShowFiles          []string
	FilesInCodeBlock   bool `default:"true"`
	MaxFileLength      int  `default:"65536"`
}

type ResultDetail struct{}

func init() {
	stage.RegisterParser(name, &ResultDetail{})
}
