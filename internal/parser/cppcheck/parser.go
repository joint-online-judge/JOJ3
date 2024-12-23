package cppcheck

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/joint-online-judge/JOJ3/internal/stage"
)

type Match struct {
	Keywords []string
	Severity []string // TODO: remove me
	Score    int
}

type Conf struct {
	Score             int
	Matches           []Match
	Stdout            string `default:"stdout"`
	Stderr            string `default:"stderr"`
	ForceQuitOnDeduct bool   `default:"false"`
}

type Record struct {
	File     string `json:"file"`
	Line     int    `json:"line"`
	Column   int    `json:"column"`
	Severity string `json:"severity"`
	Message  string `json:"message"`
	Id       string `json:"id"`
}

type CppCheck struct{}

func Parse(executorResult stage.ExecutorResult, conf Conf) stage.ParserResult {
	// stdout := executorResult.Files[conf.Stdout]
	stderr := executorResult.Files[conf.Stderr]
	records := make([]Record, 0)
	lines := strings.Split(stderr, "\n")
	for _, line := range lines {
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
	comment, score, err := GetResult(records, conf)
	if err != nil {
		return stage.ParserResult{
			Score: 0,
			Comment: fmt.Sprintf(
				"Unexpected parser error: %s.",
				err,
			),
		}
	}

	return stage.ParserResult{
		Score:   score,
		Comment: comment,
	}
}

func (*CppCheck) Run(results []stage.ExecutorResult, confAny any) (
	[]stage.ParserResult, bool, error,
) {
	conf, err := stage.DecodeConf[Conf](confAny)
	if err != nil {
		return nil, true, err
	}
	var res []stage.ParserResult
	forceQuit := false
	for _, result := range results {
		parseRes := Parse(result, *conf)
		if conf.ForceQuitOnDeduct && parseRes.Score < conf.Score {
			forceQuit = true
		}
		res = append(res, parseRes)
	}
	return res, forceQuit, nil
}
