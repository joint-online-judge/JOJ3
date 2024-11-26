package teapot

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/joint-online-judge/JOJ3/cmd/joj3/conf"
	"github.com/joint-online-judge/JOJ3/cmd/joj3/env"
)

type CheckResult struct {
	Name        string `json:"name"`
	SubmitCount int    `json:"submit_count"`
	MaxCount    int    `json:"max_count"`
	TimePeriod  int    `json:"time_period"`
}

func Check(conf *conf.Conf) (checkResults []CheckResult, err error) {
	os.Setenv("LOG_FILE_PATH", conf.Teapot.LogPath)
	os.Setenv("_TYPER_STANDARD_TRACEBACK", "1")
	if env.Attr.Actor == "" ||
		env.Attr.Repository == "" ||
		strings.Count(env.Attr.Repository, "/") != 1 {
		slog.Error("teapot env not set")
		err = fmt.Errorf("teapot env not set")
		return
	}
	repoParts := strings.Split(env.Attr.Repository, "/")
	repoName := repoParts[1]
	var formattedGroups []string
	for _, group := range conf.Teapot.Groups {
		groupConfig := fmt.Sprintf("%s=%d:%d",
			group.Name, group.MaxCount, group.TimePeriodHour)
		formattedGroups = append(formattedGroups, groupConfig)
	}
	args := []string{
		"joj3-check", conf.Teapot.EnvFilePath,
		env.Attr.Actor, conf.Teapot.GradingRepoName, repoName,
		conf.Teapot.ScoreboardPath, conf.Name,
		"--group-config", strings.Join(formattedGroups, ","),
	}
	stdoutBuf, err := runCommand(args)
	if err != nil {
		slog.Error("teapot check exec", "error", err)
		return
	}
	if json.Unmarshal(stdoutBuf.Bytes(), &checkResults) != nil {
		slog.Error("unmarshal teapot result", "error", err,
			"stdout", stdoutBuf.String())
		return
	}
	slog.Info("teapot result", "result", checkResults)
	return
}
