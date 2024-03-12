package dummy

import (
	"fmt"

	"focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/internal/stage"
)

type Conf struct {
	Score   int
	Comment string
}

type Dummy struct{}

func (*Dummy) Run(results []stage.ExecutorResult, confAny any) (
	[]stage.ParserResult, bool, error,
) {
	conf, err := stage.DecodeConf[Conf](confAny)
	if err != nil {
		return nil, true, err
	}
	var res []stage.ParserResult
	for _, result := range results {
		res = append(res, stage.ParserResult{
			Score: conf.Score,
			Comment: fmt.Sprintf(
				"%s, executor status: run time: %d ns, memory: %d bytes",
				conf.Comment, result.RunTime, result.Memory,
			),
		})
	}
	return res, false, nil
}
