package cpplint

import (
	"fmt"
	"log/slog"
	"regexp"
	"strconv"

	"github.com/joint-online-judge/JOJ3/internal/stage"
)

type Conf struct {
	Score int
}

type Cpplint struct{}

func Parse(executorResult stage.ExecutorResult, conf Conf) stage.ParserResult {
	stderr := executorResult.Files["stderr"]
	pattern := `(.+):(\d+):  (.+)  \[(.+)\] \[(\d)]\n`
	re := regexp.MustCompile(pattern)
	matches := re.FindAllStringSubmatch(stderr, -1)
	score := 0
	comment := ""
	for _, match := range matches {
		fileName := match[1]
		lineNum, err := strconv.Atoi(match[2])
		if err != nil {
			slog.Error("parse lineNum", "error", err)
			return stage.ParserResult{
				Score:   0,
				Comment: fmt.Sprintf("Unexpected parser error: %s.", err),
			}
		}
		message := match[3]
		category := match[4]
		confidence, err := strconv.Atoi(match[5])
		if err != nil {
			slog.Error("parse confidence", "error", err)
			return stage.ParserResult{
				Score:   0,
				Comment: fmt.Sprintf("Unexpected parser error: %s.", err),
			}
		}
		score -= confidence
		// TODO: add more detailed comment, just re-assemble for now
		comment += fmt.Sprintf("%s:%d:  %s  [%s] [%d]\n",
			fileName, lineNum, message, category, confidence)
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
	for _, result := range results {
		res = append(res, Parse(result, *conf))
	}
	return res, false, nil
}
