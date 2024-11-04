package cpplint

import (
	"fmt"
	"log/slog"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/joint-online-judge/JOJ3/internal/stage"
	"github.com/joint-online-judge/JOJ3/pkg/utils"
)

type Match struct {
	Keywords []string
	Score    int
}

type Conf struct {
	Score             int
	Matches           []Match
	Stdout            string `default:"stdout"`
	Stderr            string `default:"stderr"`
	ForceQuitOnDeduct bool   `default:"false"`
}

type Cpplint struct{}

func Parse(executorResult stage.ExecutorResult, conf Conf) stage.ParserResult {
	stderr := executorResult.Files[conf.Stderr]
	pattern := `(.+):(\d+):  (.+)  \[(.+)\] \[(\d)]\n`
	re := regexp.MustCompile(pattern)
	regexMatches := re.FindAllStringSubmatch(stderr, -1)
	score := conf.Score
	comment := "### Test results summary\n\n"
	categoryCount := make(map[string]int)
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
		// TODO: remove me
		if len(conf.Matches) == 0 {
			confidence, err := strconv.Atoi(regexMatch[5])
			if err != nil {
				slog.Error("parse confidence", "error", err)
				return stage.ParserResult{
					Score:   0,
					Comment: fmt.Sprintf("Unexpected parser error: %s.", err),
				}
			}
			score -= confidence
		}
		parts := strings.Split(category, "/")
		if len(parts) > 0 {
			category := parts[0]
			categoryCount[category] += 1
		}
		// TODO: remove me ends
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
	// TODO: remove me
	sortedMap := utils.SortMap(categoryCount,
		func(i, j utils.Pair[string, int]) bool {
			if i.Value == j.Value {
				return i.Key < j.Key
			}
			return i.Value > j.Value
		})
	for i, kv := range sortedMap {
		comment += fmt.Sprintf("%d. %s: %d\n", i+1, kv.Key, kv.Value)
	}
	// TODO: remove me ends
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

func (*Cpplint) Run(results []stage.ExecutorResult, confAny any) (
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
