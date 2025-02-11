package dummy

import (
	"github.com/joint-online-judge/JOJ3/internal/stage"
)

func (*Dummy) Run(results []stage.ExecutorResult, confAny any) (
	[]stage.ParserResult, bool, error,
) {
	conf, err := stage.DecodeConf[Conf](confAny)
	if err != nil {
		return nil, true, err
	}
	var res []stage.ParserResult
	for range results {
		res = append(res, stage.ParserResult{Score: conf.Score, Comment: conf.Comment})
	}
	return res, conf.ForceQuit, nil
}
