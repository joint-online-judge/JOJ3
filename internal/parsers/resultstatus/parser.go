package resultstatus

import (
	"fmt"

	"focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/internal/stage"
	"github.com/criyle/go-judge/envexec"
)

type Conf struct {
	Score   int
	Comment string
}

type ResultStatus struct{}

func (*ResultStatus) Run(results []stage.ExecutorResult, confAny any) (
	[]stage.ParserResult, bool, error,
) {
	conf, err := stage.DecodeConf[Conf](confAny)
	if err != nil {
		return nil, true, err
	}
	forceQuit := false
	var res []stage.ParserResult
	for _, result := range results {
		comment := conf.Comment
		if result.Status != stage.Status(envexec.StatusAccepted) {
			forceQuit = true
			comment = fmt.Sprintf(
				"Unexpected executor status: %s.", result.Status,
			)
		}
		res = append(res, stage.ParserResult{
			Score:   conf.Score,
			Comment: comment,
		})
	}
	return res, forceQuit, nil
}
