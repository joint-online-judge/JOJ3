package healthcheck

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"

	"github.com/joint-online-judge/JOJ3/internal/conf"
)

type CheckResult struct {
	Name        string `json:"name"`
	SubmitCount int    `json:"submit_count"`
	MaxCount    int    `json:"max_count"`
	TimePeriod  int    `json:"time_period"`
}

func runTeapot(conf *conf.Conf) (checkResults []CheckResult, err error) {
	os.Setenv("LOG_FILE_PATH", conf.Teapot.LogPath)
	var formattedGroups []string
	for _, group := range conf.Teapot.Groups {
		groupConfig := fmt.Sprintf("%s=%d:%d",
			group.Name, group.MaxCount, group.TimePeriodHour)
		formattedGroups = append(formattedGroups, groupConfig)
	}
	args := []string{
		"joj3-check-env", conf.Teapot.EnvFilePath,
		conf.Teapot.GradingRepoName,
		conf.Teapot.ScoreboardPath,
		"--group-config", strings.Join(formattedGroups, ","),
	}
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd := exec.Command("joint-teapot", args...) // #nosec G204
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf
	err = cmd.Run()
	if err != nil {
		slog.Error("teapot check exec", "error", err)
		return
	}
	if err = json.Unmarshal(stdoutBuf.Bytes(), &checkResults); err != nil {
		slog.Error("unmarshal teapot check result", "error", err,
			"stdout", stdoutBuf.String())
		return
	}
	slog.Info("teapot check result", "result", checkResults)
	return
}

func generateOutput(
	checkResults []CheckResult,
	groups []string,
) (comment string, err error) {
	if len(checkResults) == 0 {
		return
	}
	for _, checkResult := range checkResults {
		useGroup := false
		if checkResult.Name != "" {
			comment += fmt.Sprintf("keyword `%s` ", checkResult.Name)
		} else {
			useGroup = true
		}
		for _, group := range groups {
			if strings.EqualFold(group, checkResult.Name) {
				useGroup = true
				break
			}
		}
		comment += fmt.Sprintf(
			"in last %d hour(s): submit count %d, max count %d",
			checkResult.TimePeriod,
			checkResult.SubmitCount,
			checkResult.MaxCount,
		)
		if useGroup && checkResult.SubmitCount+1 > checkResult.MaxCount {
			err = fmt.Errorf(
				"keyword `%s` submit count exceeded",
				checkResult.Name,
			)
			comment += ", exceeded"
		}
		comment += "\n"
	}
	return
}

func TeapotCheck(conf *conf.Conf, groups []string) (output string, err error) {
	checkResults, err := runTeapot(conf)
	if err != nil {
		slog.Error("teapot check", "error", err)
		return
	}
	output, err = generateOutput(checkResults, groups)
	if err != nil {
		slog.Error("generate output", "error", err)
		return
	}
	return
}
