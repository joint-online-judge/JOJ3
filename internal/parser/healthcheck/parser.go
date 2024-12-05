package healthcheck

import (
	"encoding/json"
	"fmt"

	"github.com/criyle/go-judge/envexec"
	"github.com/joint-online-judge/JOJ3/internal/stage"
	"github.com/joint-online-judge/JOJ3/pkg/healthcheck"
)

type Healthcheck struct{}

type Conf struct {
	Stdout string `default:"stdout"`
	Stderr string `default:"stderr"`
}

func Parse(executorResult stage.ExecutorResult, conf Conf) (stage.ParserResult, bool) {
	stdout := executorResult.Files[conf.Stdout]
	stderr := executorResult.Files[conf.Stderr]
	if executorResult.Status != stage.Status(envexec.StatusAccepted) {
		return stage.ParserResult{
			Score: 0,
			Comment: fmt.Sprintf(
				"Unexpected executor status: `%s`\n`stderr`:\n```%s\n```\n",
				executorResult.Status, stderr,
			),
		}, true
	}
	var res healthcheck.Result
	err := json.Unmarshal([]byte(stdout), &res)
	if err != nil {
		return stage.ParserResult{
			Score: 0,
			Comment: fmt.Sprintf(
				"Failed to parse result: `%s`\n`stderr`:\n```%s\n```\n",
				err, stderr,
			),
		}, true
	}
	comment := res.Msg
	forceQuit := res.Failed
	return stage.ParserResult{
		Score:   0,
		Comment: comment,
	}, forceQuit
}

func (*Healthcheck) Run(results []stage.ExecutorResult, confAny any) (
	[]stage.ParserResult, bool, error,
) {
	conf, err := stage.DecodeConf[Conf](confAny)
	if err != nil {
		return nil, true, err
	}
	var res []stage.ParserResult
	forceQuit := false
	for _, result := range results {
		parserResult, forceQuitResult := Parse(result, *conf)
		res = append(res, parserResult)
		forceQuit = forceQuit || forceQuitResult
	}
	return res, forceQuit, nil
}
