package keyword

import (
	"fmt"
	"strings"

	"github.com/joint-online-judge/JOJ3/internal/stage"
)

type Match struct {
	Keyword       string
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
	for _, file := range conf.Files {
		content := executorResult.Files[file]
		for _, match := range conf.Matches {
			count := strings.Count(content, match.Keyword)
			if match.MaxMatchCount > 0 {
				count = min(count, match.MaxMatchCount)
			}
			if count > 0 {
				score -= count * match.Score
				comment += fmt.Sprintf(
					"Matched keyword %d time(s): %s\n",
					count, match.Keyword)
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
