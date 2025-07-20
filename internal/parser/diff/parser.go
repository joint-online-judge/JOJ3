package diff

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/joint-online-judge/JOJ3/internal/stage"
)

func (d *Diff) Run(results []stage.ExecutorResult, confAny any) (
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
		parserResult, fq, err := d.processCase(caseConf, result, conf, i)
		if err != nil {
			return nil, true, err
		}
		res = append(res, parserResult)
		if fq {
			forceQuit = true
		}
	}

	return res, forceQuit, nil
}

// processCase handles a single test case.
func (d *Diff) processCase(caseConf Case, result stage.ExecutorResult, conf *Conf, index int) (stage.ParserResult, bool, error) {
	score := 0
	comment := ""
	forceQuit := false

	if conf.FailOnNotAccepted && result.Status != stage.StatusAccepted {
		if conf.ForceQuitOnFailed {
			forceQuit = true
		}
		comment += conf.FailComment + "\n"
		comment += "Executor status not `Accepted`\n"
	} else {
		for _, output := range caseConf.Outputs {
			outputScore, outputComment, outputForceQuit, err := d.processOutput(output, result, conf, index)
			if err != nil {
				return stage.ParserResult{}, true, err
			}
			score += outputScore
			comment += outputComment
			if outputForceQuit {
				forceQuit = true
			}
		}
	}

	return stage.ParserResult{
		Score:   score,
		Comment: comment,
	}, forceQuit, nil
}

// processOutput handles a single output comparison.
func (d *Diff) processOutput(output Output, result stage.ExecutorResult, conf *Conf, index int) (int, string, bool, error) {
	answer, err := os.ReadFile(output.AnswerPath)
	if err != nil {
		return 0, "", true, err
	}
	answerStr := string(answer)
	resultStr := result.Files[output.Filename]

	isSame := stringsEqual(
		answerStr,
		resultStr,
		output.CompareSpace,
	)

	slog.Debug(
		"compare",
		"filename", output.Filename,
		"answerPath", output.AnswerPath,
		"actualLength", len(resultStr),
		"answerLength", len(answerStr),
		"index", index,
		"isSame", isSame,
	)

	if isSame {
		return output.Score, conf.PassComment + "\n", false, nil
	}

	// They are different.
	forceQuit := output.ForceQuitOnDiff
	comment := conf.FailComment + "\n"
	comment += fmt.Sprintf("Difference found in `%s`\n", output.Filename)

	if !output.AlwaysHide {
		diffComment := d.generateDiffComment(answerStr, resultStr, output)
		comment += diffComment
	} else {
		comment += "(content hidden)\n"
	}

	return 0, comment, forceQuit, nil
}

// generateDiffComment generates a diff comment for the given strings.
func (d *Diff) generateDiffComment(answerStr, resultStr string, output Output) string {
	if output.MaxDiffLength == 0 { // real default value
		output.MaxDiffLength = 4096
	}
	if output.MaxDiffLines == 0 { // real default value
		output.MaxDiffLines = 50
	}
	truncated := false

	answerLines := strings.Split(answerStr, "\n")
	resultLines := strings.Split(resultStr, "\n")
	commonPrefixLineCount := 0

	if output.HideCommonPrefix {
		n := 0
		for ; n < len(answerLines) &&
			n < len(resultLines) &&
			stringsEqual(
				answerLines[n],
				resultLines[n],
				output.CompareSpace,
			); n += 1 {
		}
		if n > 0 {
			answerLines = answerLines[n:]
			resultLines = resultLines[n:]
			commonPrefixLineCount = n
		}
	}

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
	if len(diffOutput) > output.MaxDiffLength {
		diffOutput = diffOutput[:output.MaxDiffLength]
		truncated = true
	}
	if truncated {
		diffOutput += "\n\n(truncated)"
	}
	if commonPrefixLineCount > 0 {
		diffOutput = fmt.Sprintf(
			"(%d line(s) of common prefix hidden)\n\n",
			commonPrefixLineCount,
		) + diffOutput
	}

	return fmt.Sprintf(
		"```diff\n%s\n```\n",
		diffOutput,
	)
}
