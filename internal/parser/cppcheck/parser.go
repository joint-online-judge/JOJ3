package cppcheck

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/joint-online-judge/JOJ3/internal/stage"
)

type Record struct {
	File     string `json:"file"`
	Line     int    `json:"line"`
	Column   int    `json:"column"`
	Severity string `json:"severity"`
	Message  string `json:"message"`
	ID       string `json:"id"`
}

func (*CppCheck) parse(executorResult stage.ExecutorResult, conf Conf) stage.ParserResult {
	// stdout := executorResult.Files[conf.Stdout]
	stderr := executorResult.Files[conf.Stderr]
	records := make([]Record, 0)
	lines := strings.SplitSeq(stderr, "\n")
	for line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		var record Record
		err := json.Unmarshal([]byte(line), &record)
		if err != nil {
			return stage.ParserResult{
				Score: 0,
				Comment: fmt.Sprintf(
					"Unexpected parser error: %s.",
					err,
				),
			}
		}
		records = append(records, record)
	}
	comment, score := getResult(records, conf)

	return stage.ParserResult{
		Score:   score,
		Comment: comment,
	}
}

func (p *CppCheck) Run(results []stage.ExecutorResult, confAny any) (
	[]stage.ParserResult, bool, error,
) {
	conf, err := stage.DecodeConf[Conf](confAny)
	if err != nil {
		return nil, true, err
	}
	res := make([]stage.ParserResult, 0, len(results))
	forceQuit := false
	for _, result := range results {
		parseRes := p.parse(result, *conf)
		if conf.ForceQuitOnDeduct && parseRes.Score < conf.Score {
			forceQuit = true
		}
		res = append(res, parseRes)
	}
	return res, forceQuit, nil
}
