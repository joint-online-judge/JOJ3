package clangtidy

import (
	"fmt"
	"strings"

	"github.com/criyle/go-judge/envexec"
	"github.com/joint-online-judge/JOJ3/internal/stage"
)

type Match struct {
	Keyword []string
	Score   int
}

type Conf struct {
	Score   int    `default:"100"`
	RootDir string `default:"/w"`
	Matches []Match
}

type ClangTidy struct{}

func Parse(executorResult stage.ExecutorResult, conf Conf) stage.ParserResult {
	stdout := executorResult.Files["stdout"]
	stderr := executorResult.Files["stderr"]

	lines := strings.SplitAfter(stdout, "\n")
	messages := ParseLines(lines, conf)
	formattedMessages := Format(messages)

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
	return stage.ParserResult{
		Score:   GetScore(formattedMessages, conf),
		Comment: GetComment(formattedMessages),
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
