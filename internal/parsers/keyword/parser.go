package keyword

import (
	"fmt"
	"strings"

	"focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/internal/stage"
)

type Match struct {
	Keyword string
	Score   int
}

type Conf struct {
	FullScore  int
	MinScore   int
	Files      []string
	EndOnMatch bool
	Matches    []Match
}

type Keyword struct{}

func Parse(executorResult stage.ExecutorResult, conf Conf) (
	stage.ParserResult, bool,
) {
	score := conf.FullScore
	comment := ""
	matched := false
	for _, file := range conf.Files {
		content := executorResult.Files[file]
		for _, match := range conf.Matches {
			count := strings.Count(content, match.Keyword)
			if count > 0 {
				matched = true
				score -= count * match.Score
				comment += fmt.Sprintf(
					"Matched keyword %d time(s): %s\n",
					count, match.Keyword)
			}
		}
	}
	return stage.ParserResult{
		Score:   max(score, conf.MinScore),
		Comment: comment,
	}, matched
}

func (*Keyword) Run(results []stage.ExecutorResult, confAny any) (
	[]stage.ParserResult, bool, error,
) {
	conf, err := stage.DecodeConf[Conf](confAny)
	if err != nil {
		return nil, true, err
	}
	var res []stage.ParserResult
	forceQuit := false
	for _, result := range results {
		tmp, matched := Parse(result, *conf)
		if matched && conf.EndOnMatch {
			forceQuit = true
		}
		res = append(res, tmp)
	}
	return res, forceQuit, nil
}
