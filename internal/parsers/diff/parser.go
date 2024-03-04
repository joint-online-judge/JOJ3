package diff

import (
	"fmt"
	"os"

	"focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/internal/stage"
)

type Config struct {
	Score      int
	StdoutPath string
}

type Diff struct{}

func (e *Diff) Run(result *stage.ExecutorResult, configAny any) (
	*stage.ParserResult, error,
) {
	config, err := stage.DecodeConfig[Config](configAny)
	if err != nil {
		return nil, err
	}
	score := 0
	stdout, err := os.ReadFile(config.StdoutPath)
	if err != nil {
		return nil, err
	}
	// TODO: more compare strategies
	if string(stdout) == result.Files["stdout"] {
		score = config.Score
	}
	return &stage.ParserResult{
		Score: score,
		Comment: fmt.Sprintf(
			"executor status: run time: %d ns, memory: %d bytes",
			result.RunTime, result.Memory,
		),
	}, nil
}
