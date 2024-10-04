package elf

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/criyle/go-judge/envexec"
	"github.com/joint-online-judge/JOJ3/internal/stage"
	"github.com/mitchellh/mapstructure"
)

type Conf struct {
	Score   int
	Comment string
}

type Elf struct{}

func Parse(executorResult stage.ExecutorResult, conf Conf) stage.ParserResult {
	stdout := executorResult.Files["stdout"]
	stderr := executorResult.Files["stderr"]
	if executorResult.Status != stage.Status(envexec.StatusAccepted) {
		return stage.ParserResult{
			Score: 0,
			Comment: fmt.Sprintf(
				"Unexpected executor status: %s.\nStderr: %s",
				executorResult.Status, stderr,
			),
		}
	}
	var topLevel Toplevel
	err := json.Unmarshal([]byte(stdout), &topLevel)
	if err != nil {
		return stage.ParserResult{
			Score:   0,
			Comment: fmt.Sprintf("Failed to parse result: %s", err),
		}
	}
	for _, module := range topLevel.Modules {
		for _, entry := range module.Entries {
			kind := entry[0].(string)
			report := Report{}
			err := mapstructure.Decode(entry[1], &report)
			if err != nil {
				slog.Error("elf parse", "mapstructure decode err", err)
			}
			slog.Debug("elf parse", "report file", report.File)
			slog.Debug("elf parse", "report name", report.Name)
			slog.Debug("elf parse", "report kind", kind)
			for _, caseObj := range report.Cases {
				switch kind {
				case "ParenDep":
					slog.Debug("elf parse", "binders", caseObj.Binders)
					slog.Debug("elf parse", "context", caseObj.Context)
					slog.Debug("elf parse", "depths", caseObj.Depths)
					slog.Debug("elf parse", "code", caseObj.Code)
				case "CodeLen":
					slog.Debug("elf parse", "binders", caseObj.Binders)
					slog.Debug("elf parse", "context", caseObj.Context)
					slog.Debug("elf parse", "plain", caseObj.Plain)
					slog.Debug("elf parse", "weighed", caseObj.Weighed)
					slog.Debug("elf parse", "code", caseObj.Code)
				case "OverArity":
					slog.Debug("elf parse", "binders", caseObj.Binders)
					slog.Debug("elf parse", "context", caseObj.Context)
					slog.Debug("elf parse", "detail", caseObj.Detail)
					slog.Debug("elf parse", "code", caseObj.Code)
				case "CodeDup":
					slog.Debug("elf parse", "similarity rate", caseObj.SimilarityRate)
					for _, source := range caseObj.Sources {
						slog.Debug("elf parse", "context", source.Context, "code", source.Code)
					}
				}
			}
		}
	}
	return stage.ParserResult{
		Score:   conf.Score,
		Comment: conf.Comment,
	}
}

func (*Elf) Run(results []stage.ExecutorResult, confAny any) (
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
