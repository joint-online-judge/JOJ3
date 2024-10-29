package sample

import (
	"encoding/json"
	"fmt"

	"github.com/criyle/go-judge/envexec"
	"github.com/joint-online-judge/JOJ3/internal/stage"
	"github.com/joint-online-judge/JOJ3/pkg/sample"
)

type Conf struct {
	Score   int
	Comment string
	Stdout  string `default:"stdout"`
	Stderr  string `default:"stderr"`
}

type Sample struct{}

func Parse(executorResult stage.ExecutorResult, conf Conf) stage.ParserResult {
	stdout := executorResult.Files[conf.Stdout]
	stderr := executorResult.Files[conf.Stderr]
	if executorResult.Status != stage.Status(envexec.StatusAccepted) {
		return stage.ParserResult{
			Score: 0,
			Comment: fmt.Sprintf(
				"Unexpected executor status: %s.\nStderr: %s",
				executorResult.Status, stderr,
			),
		}
	}
	var sampleResult sample.Result
	err := json.Unmarshal([]byte(stdout), &sampleResult)
	if err != nil {
		return stage.ParserResult{
			Score:   0,
			Comment: fmt.Sprintf("Failed to parse result: %s", err),
		}
	}
	return stage.ParserResult{
		Score:   sampleResult.Score + conf.Score,
		Comment: sampleResult.Comment + conf.Comment,
	}
}

func (*Sample) Run(results []stage.ExecutorResult, confAny any) (
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
