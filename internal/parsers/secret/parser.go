package secret

import (
	"focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/internal/stage"
)

type Conf struct {
	Secret string
}

type Secret struct{}

func (*Secret) Run(results []stage.ExecutorResult, confAny any) (
	[]stage.ParserResult, bool, error,
) {
	conf, err := stage.DecodeConf[Conf](confAny)
	if err != nil {
		return nil, true, err
	}
	var res []stage.ParserResult
	for range results {
		res = append(res, stage.ParserResult{Score: 0, Comment: conf.Secret})
	}
	return res, false, nil
}
