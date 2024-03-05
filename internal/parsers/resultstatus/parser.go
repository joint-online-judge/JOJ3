package resultstatus

import (
	"fmt"

	"focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/internal/stage"
	"github.com/criyle/go-judge/envexec"
)

type Config struct{}

type ResultStatus struct{}

func (*ResultStatus) Run(results []stage.ExecutorResult, configAny any) (
	[]stage.ParserResult, bool, error,
) {
	// TODO: more config options
	_, err := stage.DecodeConfig[Config](configAny)
	if err != nil {
		return nil, true, err
	}
	end := false
	var res []stage.ParserResult
	for _, result := range results {
		comment := ""
		if result.Status != stage.Status(envexec.StatusAccepted) {
			end = true
			comment = fmt.Sprintf(
				"Unexpected executor status: %s.", result.Status,
			)
		}
		res = append(res, stage.ParserResult{
			Score:   0,
			Comment: comment,
		})
	}
	return res, end, nil
}
