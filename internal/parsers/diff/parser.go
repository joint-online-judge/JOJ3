package diff

import (
	"fmt"
	"os"

	"focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/internal/stage"
)

type Config struct {
	Cases []struct {
		Score      int
		StdoutPath string
	}
}

type Diff struct{}

func (*Diff) Run(results []stage.ExecutorResult, configAny any) (
	[]stage.ParserResult, bool, error,
) {
	config, err := stage.DecodeConfig[Config](configAny)
	if err != nil {
		return nil, true, err
	}
	if len(config.Cases) != len(results) {
		return nil, true, fmt.Errorf("cases number not match")
	}
	var res []stage.ParserResult
	for i, caseConfig := range config.Cases {
		result := results[i]
		score := 0
		stdout, err := os.ReadFile(caseConfig.StdoutPath)
		if err != nil {
			return nil, true, err
		}
		// TODO: more compare strategies
		if string(stdout) == result.Files["stdout"] {
			score = caseConfig.Score
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
