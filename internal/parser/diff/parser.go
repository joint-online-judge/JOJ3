package diff

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/joint-online-judge/JOJ3/internal/stage"
)

func (*Diff) Run(results []stage.ExecutorResult, confAny any) (
	[]stage.ParserResult, bool, error,
) {
	conf, err := stage.DecodeConf[Conf](confAny)
	if err != nil {
		return nil, true, err
	}
	if len(conf.Cases) != len(results) {
		return nil, true, fmt.Errorf("cases number not match")
	}

	var res []stage.ParserResult
	forceQuit := false
	for i, caseConf := range conf.Cases {
		result := results[i]
		score := 0
		comment := ""
		if conf.FailOnNotAccepted &&
			result.Status != stage.StatusAccepted {
			if conf.ForceQuitOnFailed {
				forceQuit = true
			}
			comment += conf.FailComment
			comment += "Executor status not `Accepted`\n"
		} else {
			for _, output := range caseConf.Outputs {
				answer, err := os.ReadFile(output.AnswerPath)
				if err != nil {
					return nil, true, err
				}
				isSame := compareStrings(
					string(answer),
					result.Files[output.FileName],
					output.CompareSpace,
				)
				slog.Debug(
					"compare",
					"filename", output.FileName,
					"answerPath", output.AnswerPath,
					"actualLength", len(result.Files[output.FileName]),
					"answerLength", len(string(answer)),
					"index", i,
					"isSame", isSame,
				)
				// If no difference, assign score
				if isSame {
					score += output.Score
					comment += conf.PassComment
				} else {
					if output.ForceQuitOnDiff || conf.ForceQuitOnFailed {
						forceQuit = true
					}
					comment += conf.FailComment
					comment += fmt.Sprintf("Difference found in `%s`\n",
						output.FileName)
					if !output.AlwaysHide {
						if output.MaxDiffLength == 0 { // real default value
							output.MaxDiffLength = 2048
						}
						// Convert answer to string and split by lines
						answerStr := string(answer)
						resultStr := result.Files[output.FileName]
						if len(answerStr) > output.MaxDiffLength {
							answerStr = answerStr[:output.MaxDiffLength]
						}
						if len(resultStr) > output.MaxDiffLength {
							resultStr = resultStr[:output.MaxDiffLength]
						}
						answerLines := strings.Split(answerStr, "\n")
						resultLines := strings.Split(resultStr, "\n")
						// Generate Myers diff
						diffOps := myersDiffStr(answerLines, resultLines,
							output.CompareSpace)
						// Generate diff block with surrounding context
						diffOutput := formatDiff(
							answerLines,
							resultLines,
							diffOps,
						)
						diffOutput = strings.TrimSuffix(diffOutput, "\n  ")
						comment += fmt.Sprintf(
							"```diff\n%s\n```\n",
							diffOutput,
						)
					} else {
						comment += "(Content hidden.)\n"
					}
				}
			}
		}
		res = append(res, stage.ParserResult{
			Score:   score,
			Comment: comment,
		})
	}

	return res, forceQuit, nil
}
