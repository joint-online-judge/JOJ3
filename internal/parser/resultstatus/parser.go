package resultstatus

import (
	"fmt"

	"github.com/criyle/go-judge/envexec"
	"github.com/joint-online-judge/JOJ3/internal/stage"
)

func (*ResultStatus) Run(results []stage.ExecutorResult, confAny any) (
	[]stage.ParserResult, bool, error,
) {
	conf, err := stage.DecodeConf[Conf](confAny)
	if err != nil {
		return nil, true, err
	}
	score := conf.Score
	forceQuit := false
	var res []stage.ParserResult
	for _, result := range results {
		comment := conf.Comment
		if result.Status != stage.Status(envexec.StatusAccepted) {
			score = 0
			comment = fmt.Sprintf(
				"Unexpected executor status: `%s`.\n", result.Status,
			)
			if conf.ForceQuitOnNotAccepted {
				forceQuit = true
			}
		}
		res = append(res, stage.ParserResult{
			Score:   score,
			Comment: comment,
		})
	}
	return res, forceQuit, nil
}
