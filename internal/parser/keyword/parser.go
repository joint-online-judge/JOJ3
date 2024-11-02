package keyword

import (
	"fmt"
	"strings"

	"github.com/joint-online-judge/JOJ3/internal/stage"
)

type Match struct {
	Keywords      []string
	Keyword       string // TODO: remove me
	Score         int
	MaxMatchCount int
}

type Conf struct {
	Score             int
	FullScore         int // TODO: remove me
	Files             []string
	ForceQuitOnDeduct bool `default:"false"`
	Matches           []Match
}

type Keyword struct{}

func Parse(executorResult stage.ExecutorResult, conf Conf) stage.ParserResult {
	score := conf.Score
	comment := ""
	for _, match := range conf.Matches {
		for _, keyword := range match.Keywords {
			keywordMatchCount := 0
			for _, file := range conf.Files {
				content := executorResult.Files[file]
				keywordMatchCount += strings.Count(content, keyword)
			}
			if match.MaxMatchCount > 0 {
				keywordMatchCount = min(keywordMatchCount, match.MaxMatchCount)
			}
			if keywordMatchCount > 0 {
				score -= keywordMatchCount * match.Score
				comment += fmt.Sprintf(
					"Matched keyword %d time(s): %s\n",
					keywordMatchCount, keyword)
			}
		}
	}
	return stage.ParserResult{
		Score:   score,
		Comment: comment,
	}
}

func (*Keyword) Run(results []stage.ExecutorResult, confAny any) (
	[]stage.ParserResult, bool, error,
) {
	conf, err := stage.DecodeConf[Conf](confAny)
	if err != nil {
		return nil, true, err
	}
	// TODO: remove me on Matches.Keyword field removed
	for i := range conf.Matches {
		match := &conf.Matches[i]
		if match.Keyword != "" && len(match.Keywords) == 0 {
			match.Keywords = []string{match.Keyword}
		}
	}
	// TODO: remove me on FullScore field removed
	if conf.FullScore != 0 && conf.Score == 0 {
		conf.Score = conf.FullScore
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
