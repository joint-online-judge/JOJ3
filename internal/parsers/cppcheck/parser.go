package cppcheck

import (
	"encoding/json"
	"fmt"
	"strings"

	"focs.ji.sjtu.edu.cn/git/JOJ/JOJ3/internal/stage"
	"github.com/criyle/go-judge/envexec"
)

type CppCheck struct{}

type Match struct {
	Severity []string
	Score    int
}

type Conf struct {
	Score   int `default:"100"`
	Matches []Match
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
	// stdout := executorResult.Files["stdout"]
	stderr := executorResult.Files["stderr"]

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
		_ = json.Unmarshal([]byte(line), &record)
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
