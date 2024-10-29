package clangtidy

import (
	"fmt"
	"strings"

	"github.com/criyle/go-judge/envexec"
	"github.com/joint-online-judge/JOJ3/internal/stage"
)

type Match struct {
	Keywords []string
	Score    int
}

type Conf struct {
	Score   int    `default:"100"`
	RootDir string `default:"/w"`
	Matches []Match
	Stdout  string `default:"stdout"`
	Stderr  string `default:"stderr"`
}

type ClangTidy struct{}

func Parse(executorResult stage.ExecutorResult, conf Conf) stage.ParserResult {
	stdout := executorResult.Files[conf.Stdout]
	stderr := executorResult.Files[conf.Stderr]
	if executorResult.Status != stage.Status(envexec.StatusAccepted) {
		if !((executorResult.Status == stage.Status(envexec.StatusNonzeroExitStatus)) &&
			(executorResult.ExitStatus == 1)) {
			return stage.ParserResult{
				Score: 0,
				Comment: fmt.Sprintf(
					"Unexpected executor status: %s.\nStderr: %s",
					executorResult.Status, stderr,
				),
			}
		}
	}
	lines := strings.SplitAfter(stdout, "\n")
	messages := ParseLines(lines, conf)
	formattedMessages := Format(messages)
	score, comment := GetResult(formattedMessages, conf)
	return stage.ParserResult{
		Score:   score,
		Comment: comment,
	}
}

func (*ClangTidy) Run(results []stage.ExecutorResult, confAny any) (
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
