package cppcheck

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/criyle/go-judge/envexec"
	"github.com/joint-online-judge/JOJ3/internal/stage"
)

type CppCheck struct{}

type Match struct {
	Severity []string
	Score    int
}

type Conf struct {
	Score   int
	Matches []Match
	Stdout  string `default:"stdout"`
	Stderr  string `default:"stderr"`
}

type Record struct {
	File     string `json:"file"`
	Line     int    `json:"line"`
	Column   int    `json:"column"`
	Severity string `json:"severity"`
	Message  string `json:"message"`
	Id       string `json:"id"`
}

func Parse(executorResult stage.ExecutorResult, conf Conf) stage.ParserResult {
	// stdout := executorResult.Files[conf.Stdout]
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
	for _, result := range results {
		res = append(res, Parse(result, *conf))
	}
	return res, false, nil
}
