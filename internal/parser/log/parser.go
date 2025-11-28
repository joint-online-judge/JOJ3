package log

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/joint-online-judge/JOJ3/internal/stage"
)

func (*Log) parse(executorResult stage.ExecutorResult, conf Conf) stage.ParserResult {
	content := executorResult.Files[conf.Filename]
	var data map[string]any
	contentBytes := []byte(content)
	if json.Valid(contentBytes) {
		err := json.Unmarshal(contentBytes, &data)
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
		slog.Default().Log(
			context.Background(),
			slog.Level(conf.Level),
			conf.Msg,
			args...,
		)
	} else {
		slog.Default().Log(
			context.Background(),
			slog.Level(conf.Level),
			conf.Msg,
			"content",
			content,
		)
	}
	return stage.ParserResult{
		Score:   0,
		Comment: "",
	}
}

func (p *Log) Run(results []stage.ExecutorResult, confAny any) (
	[]stage.ParserResult, bool, error,
) {
	conf, err := stage.DecodeConf[Conf](confAny)
	if err != nil {
		return nil, true, err
	}
	res := make([]stage.ParserResult, 0, len(results))
	for _, result := range results {
		res = append(res, p.parse(result, *conf))
	}
	return res, false, nil
}
