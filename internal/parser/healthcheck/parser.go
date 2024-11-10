package healthcheck

import (
	"fmt"

	"github.com/criyle/go-judge/envexec"
	"github.com/joint-online-judge/JOJ3/internal/stage"
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
				"Unexpected executor status: `%s`\n`stdout`:\n```%s\n```\n`stderr`:\n```%s\n```",
				executorResult.Status, stdout, stderr,
			),
		}, true
	}
	return stage.ParserResult{
		Score:   0,
		Comment: stdout,
	}, stdout != ""
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
