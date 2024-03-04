package dummy

import (
	"fmt"
	"log/slog"

	"focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/internal/stage"
	"github.com/mitchellh/mapstructure"
)

type Config struct {
	Score   int
	Comment string
}

type Dummy struct{}

func (e *Dummy) Run(result *stage.Result, configAny any) (
	*stage.ParserResult, error,
) {
	var config Config
	err := mapstructure.Decode(configAny, &config)
	if err != nil {
		slog.Error("failed to decode config", "err", err)
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
