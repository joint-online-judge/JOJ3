package resultstatus

import (
	"fmt"

	"github.com/joint-online-judge/JOJ3/internal/stage"
)

// statusCodes uses the abbreviations accepted by result-status.whitelisted-codes.
var statusCodes = map[string]stage.Status{
	"MLE": stage.StatusMemoryLimitExceeded,
	"TLE": stage.StatusTimeLimitExceeded,
	"OLE": stage.StatusOutputLimitExceeded,
	"FE":  stage.StatusFileError,
	"NZS": stage.StatusNonzeroExitStatus,
	"SIG": stage.StatusSignalled,
	"DJS": stage.StatusDangerousSyscall,
	"IE":  stage.StatusInternalError,
}

func (*ResultStatus) Run(results []stage.ExecutorResult, confAny any) (
	[]stage.ParserResult, bool, error,
) {
	conf, err := stage.DecodeConf[Conf](confAny)
	if err != nil {
		return nil, true, err
	}
	whitelisted := make(map[stage.Status]bool, len(conf.WhitelistedCodes))
	for _, code := range conf.WhitelistedCodes {
		status, ok := statusCodes[code]
		if !ok {
			return nil, true, fmt.Errorf("invalid result-status whitelisted code: %q", code)
		}
		whitelisted[status] = true
	}

	forceQuit := false
	res := make([]stage.ParserResult, 0, len(results))
	for _, result := range results {
		score := conf.Score
		comment := conf.Comment
		if result.Status != stage.StatusAccepted {
			score = 0
			comment = fmt.Sprintf(
				"Unexpected executor status: `%s`.\n", result.Status,
			)
			if conf.ForceQuitOnNotAccepted && !whitelisted[result.Status] {
				forceQuit = true
			}
		}
		res = append(res, stage.ParserResult{
			Score:   score,
			Comment: comment,
		})
	}
	return res, forceQuit, nil
}
