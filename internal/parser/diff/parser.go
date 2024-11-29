package diff

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/criyle/go-judge/envexec"
	"github.com/joint-online-judge/JOJ3/internal/stage"
)

// operation represents the type of edit operation.
type operation uint

const (
	INSERT operation = iota + 1
	DELETE
	MOVE
)

type Conf struct {
	PassComment       string `default:"🥳Passed!\n"`
	FailComment       string `default:"🧐Failed...\n"`
	FailOnNotAccepted bool   `default:"true"`
	ForceQuitOnFailed bool   `default:"false"`
	Cases             []struct {
		Outputs []struct {
			Score           int
			FileName        string
			AnswerPath      string
			CompareSpace    bool
			AlwaysHide      bool
			ForceQuitOnDiff bool
			MaxDiffLength   int `default:"2048"` // just for reference
		}
	}
}

type Diff struct{}

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
			result.Status != stage.Status(envexec.StatusAccepted) {
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
				slog.Debug("compare", "filename", output.FileName,
					"answer path", output.AnswerPath,
					"actual length", len(result.Files[output.FileName]),
					"answer length", len(string(answer)))
				// If no difference, assign score
				if compareStrings(string(answer), result.Files[output.FileName],
					output.CompareSpace) {
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
						// Convert answer to string and split by lines
						stdoutLines := strings.Split(string(answer), "\n")
						resultLines := strings.Split(
							result.Files[output.FileName], "\n")

						// Generate Myers diff
						diffOps := myersDiff(stdoutLines, resultLines,
							output.CompareSpace)
						if output.MaxDiffLength == 0 { // real default value
							output.MaxDiffLength = 2048
						}
						// Generate diff block with surrounding context
						diffOutput := generateDiffWithContext(
							stdoutLines,
							resultLines,
							diffOps,
							output.MaxDiffLength,
						)
						diffOutput = strings.TrimSuffix(diffOutput, "\n  \n")
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
