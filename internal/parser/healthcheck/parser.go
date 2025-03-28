package healthcheck

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/joint-online-judge/JOJ3/internal/stage"
	"github.com/joint-online-judge/JOJ3/pkg/healthcheck"
)

func (*Healthcheck) parse(executorResult stage.ExecutorResult, conf Conf) (stage.ParserResult, bool) {
	stdout := executorResult.Files[conf.Stdout]
	stderr := executorResult.Files[conf.Stderr]
	slog.Debug("healthcheck files", "stdout", stdout, "stderr", stderr)
	if executorResult.Status != stage.StatusAccepted {
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
	slog.Debug("healthcheck result", "res", res)
	comment := res.Msg
	forceQuit := res.Failed
	return stage.ParserResult{
		Score:   0,
		Comment: comment,
	}, forceQuit
}

func (p *Healthcheck) Run(results []stage.ExecutorResult, confAny any) (
	[]stage.ParserResult, bool, error,
) {
	conf, err := stage.DecodeConf[Conf](confAny)
	if err != nil {
		return nil, true, err
	}
	res := make([]stage.ParserResult, 0, len(results))
	forceQuit := false
	for _, result := range results {
		parserResult, forceQuitResult := p.parse(result, *conf)
		res = append(res, parserResult)
		forceQuit = forceQuit || forceQuitResult
	}
	return res, forceQuit, nil
}
