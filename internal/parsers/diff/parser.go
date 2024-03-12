package diff

import (
	"fmt"
	"os"

	"focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/internal/stage"
)

type Conf struct {
	Cases []struct {
		Score      int
		StdoutPath string
	}
}

type Diff struct{}

func (*Diff) Run(results []stage.ExecutorResult, confAny any) (
	[]stage.ParserResult, bool, error,
) {
	conf, err := stage.DecodeConf[Conf](confAny)
	if err != nil {
		return nil, true, err
	}
	if len(conf.Cases) != len(results) {
		return nil, true, fmt.Errorf("cases number not match")
	}
	var res []stage.ParserResult
	for i, caseConf := range conf.Cases {
		result := results[i]
		score := 0
		stdout, err := os.ReadFile(caseConf.StdoutPath)
		if err != nil {
			return nil, true, err
		}
		// TODO: more compare strategies
		if string(stdout) == result.Files["stdout"] {
			score = caseConf.Score
		}
		res = append(res, stage.ParserResult{
			Score: score,
			Comment: fmt.Sprintf(
				"executor status: run time: %d ns, memory: %d bytes",
				result.RunTime, result.Memory,
			),
		})
	}
	return res, false, nil
}
