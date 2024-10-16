package resultdetail

import (
	"fmt"

	"github.com/joint-online-judge/JOJ3/internal/stage"
)

type Conf struct {
	Score          int
	ShowExitStatus bool `default:"false"`
	ShowError      bool `default:"false"`
	ShowTime       bool `default:"true"`
	ShowMemory     bool `default:"true"`
	ShowRunTime    bool `default:"false"`
	ShowFiles      []string
}

type ResultDetail struct{}

func (*ResultDetail) Run(results []stage.ExecutorResult, confAny any) (
	[]stage.ParserResult, bool, error,
) {
	conf, err := stage.DecodeConf[Conf](confAny)
	if err != nil {
		return nil, true, err
	}
	forceQuit := false
	var res []stage.ParserResult
	for _, result := range results {
		comment := ""
		if conf.ShowExitStatus {
			comment += fmt.Sprintf("Exit Status: `%d`\n", result.ExitStatus)
		}
		if conf.ShowError {
			if result.Error == "" {
				result.Error = "nil"
			}
			comment += fmt.Sprintf("Error: `%s`\n", result.Error)
		}
		if conf.ShowTime {
			comment += fmt.Sprintf("Time: `%d ms`\n", result.Time/1e6)
		}
		if conf.ShowMemory {
			comment += fmt.Sprintf("Memory: `%.2f MiB`\n",
				float64(result.Memory)/(1024*1024))
		}
		if conf.ShowRunTime {
			comment += fmt.Sprintf("RunTime: `%d ms`\n", result.RunTime/1e6)
		}
		for _, file := range conf.ShowFiles {
			content, ok := result.Files[file]
			comment += fmt.Sprintf("File `%s`:\n", file)
			if ok {
				comment += fmt.Sprintf("```\n%s\n```\n", content)
			} else {
				comment += "Not found.\n"
			}
		}
		res = append(res, stage.ParserResult{
			Score:   conf.Score,
			Comment: comment,
		})
	}
	return res, forceQuit, nil
}
