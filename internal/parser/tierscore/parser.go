package tierscore

import (
	"fmt"
	"strings"

	"github.com/joint-online-judge/JOJ3/internal/stage"
)

func (*TierScore) Run(results []stage.ExecutorResult, confAny any) (
	[]stage.ParserResult, bool, error,
) {
	conf, err := stage.DecodeConf[Conf](confAny)
	if err != nil {
		return nil, true, err
	}

	var res []stage.ParserResult
	forceQuit := false

	for _, result := range results {
		totalScore := 0
		var commentBuilder strings.Builder

		for i, tier := range conf.Tiers {
			conditionsMet := true
			var conditionDesc []string

			if tier.TimeLessThan > 0 {
				if result.Time < tier.TimeLessThan {
					conditionDesc = append(
						conditionDesc,
						fmt.Sprintf(
							"Time < `%d ms`",
							tier.TimeLessThan/1e6,
						),
					)
				} else {
					conditionsMet = false
				}
			}

			if tier.MemoryLessThan > 0 {
				if result.Memory < tier.MemoryLessThan {
					conditionDesc = append(
						conditionDesc,
						fmt.Sprintf(
							"Memory < `%.2f MiB`",
							float64(tier.MemoryLessThan)/(1024*1024),
						),
					)
				} else {
					conditionsMet = false
				}
			}

			if conditionsMet {
				totalScore += tier.Score
				fmt.Fprintf(
					&commentBuilder,
					"Tier %d: +%d (meets %s)\n",
					i,
					tier.Score,
					strings.Join(conditionDesc, " and "),
				)
			} else {
				var required []string
				if tier.TimeLessThan > 0 {
					required = append(
						required,
						fmt.Sprintf(
							"Time < `%d ms`",
							tier.TimeLessThan/1e6,
						),
					)
				}
				if tier.MemoryLessThan > 0 {
					required = append(
						required,
						fmt.Sprintf(
							"Memory < `%.2f MiB`",
							float64(tier.MemoryLessThan)/(1024*1024),
						),
					)
				}
				fmt.Fprintf(
					&commentBuilder,
					"Tier %d: +0 (requires %s)\n",
					i+1,
					strings.Join(required, " and "),
				)
			}
		}
		fmt.Fprintf(&commentBuilder, "Final score: %d", totalScore)
		res = append(res, stage.ParserResult{
			Score:   totalScore,
			Comment: commentBuilder.String(),
		})
	}

	return res, forceQuit, nil
}
