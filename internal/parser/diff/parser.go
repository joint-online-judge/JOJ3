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

	res := make([]stage.ParserResult, 0, len(conf.Cases))
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
				isSame := stringsEqual(
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
						if output.MaxDiffLines == 0 { // real default value
							output.MaxDiffLines = 50
						}
						// Convert answer to string and split by lines
						truncated := false
						answerStr := string(answer)
						resultStr := result.Files[output.FileName]
						if len(answerStr) > output.MaxDiffLength {
							answerStr = answerStr[:output.MaxDiffLength]
							truncated = true
						}
						if len(resultStr) > output.MaxDiffLength {
							resultStr = resultStr[:output.MaxDiffLength]
							truncated = true
						}
						answerLines := strings.Split(answerStr, "\n")
						resultLines := strings.Split(resultStr, "\n")
						if len(answerLines) > output.MaxDiffLines {
							answerLines = answerLines[:output.MaxDiffLines]
							truncated = true
						}
						if len(resultLines) > output.MaxDiffLines {
							resultLines = resultLines[:output.MaxDiffLines]
							truncated = true
						}
						diffs := patienceDiff(
							answerLines,
							resultLines,
							func(a, b string) bool {
								return stringsEqual(a, b, output.CompareSpace)
							})
						diffOutput := diffText(diffs)
						diffOutput = strings.TrimSuffix(diffOutput, "\n  ")
						if truncated {
							diffOutput += "\n\n(truncated)"
						}
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
