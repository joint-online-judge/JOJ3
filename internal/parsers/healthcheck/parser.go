package healthcheck

import (
	"fmt"

	"focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/internal/stage"
	"github.com/criyle/go-judge/envexec"
)

type Healthcheck struct{}

func Parse(executorResult stage.ExecutorResult) (stage.ParserResult, bool) {
	stdout := executorResult.Files["stdout"]
	stderr := executorResult.Files["stderr"]
	if executorResult.Status != stage.Status(envexec.StatusAccepted) {
		return stage.ParserResult{
			Score: 0,
			Comment: fmt.Sprintf(
				"Unexpected executor status: %s.\nStdout: %s\nStderr: %s",
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
	var res []stage.ParserResult
	forceQuit := false
	for _, result := range results {
		parserResult, forceQuitResult := Parse(result)
		res = append(res, parserResult)
		forceQuit = forceQuit || forceQuitResult
	}
	return res, forceQuit, nil
}
