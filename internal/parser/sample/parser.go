package sample

import (
	"encoding/json"
	"fmt"

	"github.com/joint-online-judge/JOJ3/internal/stage"
	"github.com/joint-online-judge/JOJ3/pkg/sample"
)

func (*Sample) parse(executorResult stage.ExecutorResult, conf Conf) stage.ParserResult {
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

func (p *Sample) Run(results []stage.ExecutorResult, confAny any) (
	[]stage.ParserResult, bool, error,
) {
	conf, err := stage.DecodeConf[Conf](confAny)
	if err != nil {
		return nil, true, err
	}
	res := make([]stage.ParserResult, 0, len(results))
	for _, result := range results {
		res = append(res, p.parse(result, *conf))
	}
	return res, false, nil
}
