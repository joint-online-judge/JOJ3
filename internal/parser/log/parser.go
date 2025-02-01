package log

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/joint-online-judge/JOJ3/internal/stage"
)

type Conf struct {
	FileName string `default:"stdout"`
	Msg      string `default:"log msg"`
}

type Log struct{}

func Parse(executorResult stage.ExecutorResult, conf Conf) stage.ParserResult {
	content := executorResult.Files[conf.FileName]
	var data map[string]any
	err := json.Unmarshal([]byte(content), &data)
	if err != nil {
		slog.Error(conf.Msg, "error", err)
		return stage.ParserResult{
			Score:   0,
			Comment: fmt.Sprintf("Failed to parse content: %s", err),
		}
	}
	args := make([]any, 0, len(data)*2)
	for key, value := range data {
		args = append(args, key, value)
	}
	slog.Info(conf.Msg, args...)
	return stage.ParserResult{
		Score:   0,
		Comment: "",
	}
}

func (*Log) Run(results []stage.ExecutorResult, confAny any) (
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
