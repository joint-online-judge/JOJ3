package clangtidy

import (
	"strings"

	"github.com/joint-online-judge/JOJ3/internal/stage"
)

func Parse(executorResult stage.ExecutorResult, conf Conf) stage.ParserResult {
	stdout := executorResult.Files[conf.Stdout]
	// stderr := executorResult.Files[conf.Stderr]
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
