package elf

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/joint-online-judge/JOJ3/internal/stage"
	"github.com/mitchellh/mapstructure"
)

func (p *Elf) parse(executorResult stage.ExecutorResult, conf Conf) stage.ParserResult {
	stdout := executorResult.Files[conf.Stdout]
	// stderr := executorResult.Files[conf.Stderr]
	var topLevel Toplevel
	err := json.Unmarshal([]byte(stdout), &topLevel)
	if err != nil {
		return stage.ParserResult{
			Score: 0,
			Comment: fmt.Sprintf(
				"Unexpected parser error: %s.",
				err,
			),
		}
	}
	score := conf.Score
	comment := ""
	for _, module := range topLevel.Modules {
		for _, entry := range module.Entries {
			kind := entry[0].(string)
			report := Report{}
			err := mapstructure.Decode(entry[1], &report)
			if err != nil {
				slog.Error("elf parse", "mapstructure decode err", err)
			}
			comment += fmt.Sprintf("### [%s] %s\n", report.File, report.Name)
			for _, caseObj := range report.Cases {
				for _, match := range conf.Matches {
					for _, keyword := range match.Keywords {
						if strings.Contains(kind, keyword) {
							score += -match.Score
						}
					}
				}
				switch kind {
				case "ParenDep":
					// "<binders>:\n<context> below reaches a parentheses depths of <depths>:\n<code>"
					comment += fmt.Sprintf(
						"%s:\n%s below reaches a parentheses depths of %d:\n"+
							"```%s\n```\n",
						caseObj.Binders,
						caseObj.Context,
						caseObj.Depths,
						caseObj.Code,
					)
				case "CodeLen":
					// "<binders>:\n<context> below excceeds a code length upper bound with <plain> (weighed: <weighed>):\n<code>"
					comment += fmt.Sprintf(
						"%s:\n%s below excceeds a code length "+
							"upper bound with %d (weighed: %f):\n"+
							"```%s\n```\n",
						caseObj.Binders,
						caseObj.Context,
						caseObj.Plain,
						caseObj.Weighed,
						caseObj.Code,
					)
				case "OverArity":
					// "<binders>:\n<context> below hits <detail>:\n<code>"
					comment += fmt.Sprintf(
						"%s:\n%s below hits %s:\n```%s\n```\n",
						caseObj.Binders,
						caseObj.Context,
						caseObj.Detail,
						caseObj.Code,
					)
				case "CodeDup":
					if len(caseObj.Sources) != 2 {
						slog.Error("elf parse", "code dup sources length", len(caseObj.Sources))
					}
					context0 := caseObj.Sources[0].Context
					code0 := caseObj.Sources[0].Code
					context1 := caseObj.Sources[1].Context
					code1 := caseObj.Sources[1].Code
					// "The code below has a similarity rate of <similarity_rate>:\n- <context1>:\n\t<code1>\n- <context2>:\n\t<code2>"
					comment += fmt.Sprintf(
						"The code below has a similarity rate of %f:\n"+
							"- %s:\n```%s\n```\n"+
							"- %s:\n```%s\n```\n",
						caseObj.SimilarityRate,
						context0,
						code0,
						context1,
						code1,
					)
				}
			}
		}
	}
	return stage.ParserResult{
		Score:   score,
		Comment: comment,
	}
}

func (p *Elf) Run(results []stage.ExecutorResult, confAny any) (
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
