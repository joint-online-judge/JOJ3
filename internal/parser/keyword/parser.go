package keyword

import (
	"fmt"
	"strings"

	"github.com/joint-online-judge/JOJ3/internal/stage"
)

type Match struct {
	Keyword string
	Score   int
}

type Conf struct {
	Score            int
	FullScore        int // TODO: remove me
	MinScore         int
	Files            []string
	ForceQuitOnMatch bool
	Matches          []Match
}

type Keyword struct{}

func Parse(executorResult stage.ExecutorResult, conf Conf) (
	stage.ParserResult, bool,
) {
	score := conf.Score
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
	// TODO: remove me on FullScore field removed
	if conf.FullScore != 0 && conf.Score == 0 {
		conf.Score = conf.FullScore
	}
	var res []stage.ParserResult
	forceQuit := false
	for _, result := range results {
		tmp, matched := Parse(result, *conf)
		if matched && conf.ForceQuitOnMatch {
			forceQuit = true
		}
		res = append(res, tmp)
	}
	return res, forceQuit, nil
}
