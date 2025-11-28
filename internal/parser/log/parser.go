package log

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/joint-online-judge/JOJ3/internal/stage"
)

func (*Log) parse(executorResult stage.ExecutorResult, conf Conf) stage.ParserResult {
	content, ok := executorResult.Files[conf.Filename]
	if !ok {
		slog.Error("file not found for log parser", "filename", conf.Filename)
		return stage.ParserResult{
			Score:   0,
			Comment: fmt.Sprintf("log parser: file %s not found", conf.Filename),
		}
	}
	contentBytes := []byte(content)
	var data map[string]any
	if err := json.Unmarshal(contentBytes, &data); err != nil {
		// Not a valid json or failed to unmarshal, log as raw string line by line.
		for line := range strings.SplitSeq(content, "\n") {
			if strings.TrimSpace(line) != "" {
				slog.Default().Log(
					context.Background(),
					slog.Level(conf.Level),
					conf.Msg,
					"line",
					line,
				)
			}
		}
	} else {
		// Valid json, log as key-value pairs.
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
