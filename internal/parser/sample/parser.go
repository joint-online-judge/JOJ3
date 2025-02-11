package sample

import (
	"encoding/json"
	"fmt"

	"github.com/joint-online-judge/JOJ3/internal/stage"
	"github.com/joint-online-judge/JOJ3/pkg/sample"
)

func Parse(executorResult stage.ExecutorResult, conf Conf) stage.ParserResult {
	stdout := executorResult.Files[conf.Stdout]
	// stderr := executorResult.Files[conf.Stderr]
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
