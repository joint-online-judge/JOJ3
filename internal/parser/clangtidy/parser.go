package clangtidy

import (
	"strings"

	"github.com/joint-online-judge/JOJ3/internal/stage"
)

func (*ClangTidy) parse(executorResult stage.ExecutorResult, conf Conf) stage.ParserResult {
	stdout := executorResult.Files[conf.Stdout]
	// stderr := executorResult.Files[conf.Stderr]
	lines := strings.SplitAfter(stdout, "\n")
	messages := parseLines(lines, conf)
	formattedMessages := format(messages)
	score, comment := getResult(formattedMessages, conf)
	return stage.ParserResult{
		Score:   score,
		Comment: comment,
	}
}

func (p *ClangTidy) Run(results []stage.ExecutorResult, confAny any) (
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
