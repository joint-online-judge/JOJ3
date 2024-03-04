package dummy

import (
	"fmt"

	"focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/internal/stage"
)

type Config struct {
	Score   int
	Comment string
}

type Dummy struct{}

func (e *Dummy) Run(result *stage.ExecutorResult, configAny any) (
	*stage.ParserResult, error,
) {
	config, err := stage.DecodeConfig[Config](configAny)
	if err != nil {
		return nil, err
	}
	return &stage.ParserResult{
		Score: config.Score,
		Comment: fmt.Sprintf(
			"%s, executor status: run time: %d ns, memory: %d bytes",
			config.Comment, result.RunTime, result.Memory,
		),
	}, nil
}
