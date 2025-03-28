package cpplint

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/joint-online-judge/JOJ3/internal/stage"
)

func (*Cpplint) parse(executorResult stage.ExecutorResult, conf Conf) stage.ParserResult {
	stderr := executorResult.Files[conf.Stderr]
	pattern := `(.+):(\d+):  (.+)  \[(.+)\] \[(\d)]\n`
	re := regexp.MustCompile(pattern)
	regexMatches := re.FindAllStringSubmatch(stderr, -1)
	score := conf.Score
	comment := "### Test results summary\n\n"
	matchCount := make(map[string]int)
	scoreChange := make(map[string]int)
	for _, regexMatch := range regexMatches {
		// fileName := regexMatch[1]
		// lineNum, err := strconv.Atoi(regexMatch[2])
		// if err != nil {
		// 	slog.Error("parse lineNum", "error", err)
		// 	return stage.ParserResult{
		// 		Score:   0,
		// 		Comment: fmt.Sprintf("Unexpected parser error: %s.", err),
		// 	}
		// }
		// message := regexMatch[3]
		category := regexMatch[4]
		for _, match := range conf.Matches {
			for _, keyword := range match.Keywords {
				if strings.Contains(category, keyword) {
					matchCount[keyword] += 1
					scoreChange[keyword] += -match.Score
					score += -match.Score
				}
			}
		}
	}
	type Result struct {
		Keyword     string
		Count       int
		ScoreChange int
	}
	results := make([]Result, 0, len(matchCount))
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

func (p *Cpplint) Run(results []stage.ExecutorResult, confAny any) (
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
