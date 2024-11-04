package keyword

import (
	"fmt"
	"sort"
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
	matchCount := make(map[string]int)
	scoreChange := make(map[string]int)
	for _, match := range conf.Matches {
		for _, keyword := range match.Keywords {
			for _, file := range conf.Files {
				content := executorResult.Files[file]
				matchCount[keyword] += strings.Count(content, keyword)
			}
			if match.MaxMatchCount > 0 {
				matchCount[keyword] = min(
					matchCount[keyword], match.MaxMatchCount)
			}
			score += -match.Score * matchCount[keyword]
			scoreChange[keyword] = -match.Score * matchCount[keyword]
		}
	}
	type Result struct {
		Keyword     string
		Count       int
		ScoreChange int
	}
	var results []Result
	for keyword, count := range matchCount {
		results = append(results, Result{
			Keyword:     keyword,
			Count:       count,
			ScoreChange: scoreChange[keyword],
		})
	}
	sort.Slice(results, func(i, j int) bool {
		if results[i].ScoreChange != results[j].ScoreChange {
			return results[i].ScoreChange < results[j].ScoreChange
		}
		if results[i].Count != results[j].Count {
			return results[i].Count > results[j].Count
		}
		return results[i].Keyword < results[j].Keyword
	})
	for i, result := range results {
		comment += fmt.Sprintf("%d. `%s`: %d occurrence(s), %d point(s)\n",
			i+1, result.Keyword, result.Count, result.ScoreChange)
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
