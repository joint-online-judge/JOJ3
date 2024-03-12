package dummy

import (
	"encoding/json"
	"fmt"

	"focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/internal/stage"
	"focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/pkg/dummy"
	"github.com/criyle/go-judge/envexec"
)

type Conf struct {
	Score   int
	Comment string
}

type Dummy struct{}

func Parse(executorResult stage.ExecutorResult, conf Conf) stage.ParserResult {
	stdout := executorResult.Files["stdout"]
	stderr := executorResult.Files["stderr"]
	if executorResult.Status != stage.Status(envexec.StatusAccepted) {
		return stage.ParserResult{
			Score: 0,
			Comment: fmt.Sprintf(
				"Unexpected executor status: %s.\nStderr: %s",
				executorResult.Status, stderr,
			),
		}
	}
	var dummyResult dummy.Result
	err := json.Unmarshal([]byte(stdout), &dummyResult)
	if err != nil {
		return stage.ParserResult{
			Score:   0,
			Comment: fmt.Sprintf("Failed to parse result: %s", err),
		}
	}
	return stage.ParserResult{
		Score:   dummyResult.Score + conf.Score,
		Comment: dummyResult.Comment + conf.Comment,
	}
}

func (*Dummy) Run(results []stage.ExecutorResult, confAny any) (
	[]stage.ParserResult, bool, error,
) {
	conf, err := stage.DecodeConf[Conf](confAny)
	if err != nil {
		return nil, true, err
	}
	var res []stage.ParserResult
	for _, result := range results {
		res = append(res, Parse(result, *conf))
	}
	return res, false, nil
}
