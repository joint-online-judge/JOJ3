package debug

import (
	"log/slog"

	"github.com/joint-online-judge/JOJ3/internal/stage"
)

func (*Debug) parse(executorResult stage.ExecutorResult, _ Conf) stage.ParserResult {
	slog.Debug("debug parser", "executorResult", executorResult)
	for name, content := range executorResult.Files {
		slog.Debug("debug parser file", "name", name, "content", content)
	}
	return stage.ParserResult{
		Score:   0,
		Comment: "",
	}
}

func (p *Debug) Run(results []stage.ExecutorResult, confAny any) (
	[]stage.ParserResult, bool, error,
) {
	conf, err := stage.DecodeConf[Conf](confAny)
	if err != nil {
		return nil, true, err
	}
	var res []stage.ParserResult
	for _, result := range results {
		res = append(res, p.parse(result, *conf))
	}
	return res, false, nil
}
