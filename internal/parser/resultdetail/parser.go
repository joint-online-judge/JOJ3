package resultdetail

import (
	"fmt"

	"github.com/joint-online-judge/JOJ3/internal/stage"
)

func (*ResultDetail) Run(results []stage.ExecutorResult, confAny any) (
	[]stage.ParserResult, bool, error,
) {
	conf, err := stage.DecodeConf[Conf](confAny)
	if err != nil {
		return nil, true, err
	}
	forceQuit := false
	res := make([]stage.ParserResult, 0, len(results))
	for _, result := range results {
		comment := ""
		if conf.ShowExecutorStatus {
			comment += fmt.Sprintf("Executor Status: `%s`\n", result.Status.String())
		}
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
			comment += fmt.Sprintf("Wall-clock Time: `%d ms`\n", result.RunTime/1e6)
		}
		if conf.ShowProcPeak {
			comment += fmt.Sprintf("ProcPeak: `%d`\n", result.ProcPeak)
		}
		for _, file := range conf.ShowFiles {
			content, ok := result.Files[file]
			comment += fmt.Sprintf("File `%s`:\n", file)
			if ok {
				if conf.MaxFileLength > 0 && len(content) > conf.MaxFileLength {
					content = content[:conf.MaxFileLength] + "\n\n(truncated)"
				}
				if conf.FilesInCodeBlock {
					comment += fmt.Sprintf("`````````\n%s\n`````````\n", content)
				} else {
					comment += fmt.Sprintf("%s\n", content)
				}
			} else {
				comment += "Not found\n"
			}
		}
		res = append(res, stage.ParserResult{
			Score:   conf.Score,
			Comment: comment,
		})
	}
	return res, forceQuit, nil
}
