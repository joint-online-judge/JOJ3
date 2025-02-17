package keyword

import (
	"fmt"
	"sort"
	"strings"

	"github.com/joint-online-judge/JOJ3/internal/stage"
)

func (*Keyword) parse(executorResult stage.ExecutorResult, conf Conf) stage.ParserResult {
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

func (k *Keyword) Run(results []stage.ExecutorResult, confAny any) (
	[]stage.ParserResult, bool, error,
) {
	conf, err := stage.DecodeConf[Conf](confAny)
	if err != nil {
		return nil, true, err
	}
	var res []stage.ParserResult
	forceQuit := false
	for _, result := range results {
		parseRes := k.parse(result, *conf)
		if conf.ForceQuitOnDeduct && parseRes.Score < conf.Score {
			forceQuit = true
		}
		res = append(res, parseRes)
	}
	return res, forceQuit, nil
}
