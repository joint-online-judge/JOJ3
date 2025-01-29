package teapot

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/joint-online-judge/JOJ3/cmd/joj3/env"
	"github.com/joint-online-judge/JOJ3/internal/conf"
)

type RunResult struct {
	Issue  int    `json:"issue"`
	Action int    `json:"action"`
	Sha    string `json:"sha"`
}

func Run(conf *conf.Conf) (
	runResult RunResult, err error,
) {
	os.Setenv("LOG_FILE_PATH", conf.Teapot.LogPath)
	os.Setenv("_TYPER_STANDARD_TRACEBACK", "1")
	if env.Attr.Actor == "" ||
		env.Attr.Repository == "" ||
		strings.Count(env.Attr.Repository, "/") != 1 ||
		env.Attr.RunNumber == "" {
		slog.Error("teapot env not set")
		err = fmt.Errorf("teapot env not set")
		return
	}
	repoParts := strings.Split(env.Attr.Repository, "/")
	repoName := repoParts[1]
	skipIssueArg := "--no-skip-result-issue"
	if conf.Teapot.SkipIssue {
		skipIssueArg = "--skip-result-issue"
	}
	skipScoreboardArg := "--no-skip-scoreboard"
	if conf.Teapot.SkipScoreboard {
		skipScoreboardArg = "--skip-scoreboard"
	}
	skipFailedTableArg := "--no-skip-failed-table"
	if conf.Teapot.SkipFailedTable {
		skipFailedTableArg = "--skip-failed-table"
	}
	submitterInIssueTitleArg := "--no-submitter-in-issue-title"
	if conf.Teapot.SubmitterInIssueTitle {
		submitterInIssueTitleArg = "--submitter-in-issue-title"
	}
	args := []string{
		"joj3-all", conf.Teapot.EnvFilePath, conf.Stage.OutputPath,
		env.Attr.Actor, conf.Teapot.GradingRepoName, repoName,
		env.Attr.RunNumber, conf.Teapot.ScoreboardPath,
		conf.Teapot.FailedTablePath,
		conf.Name, env.Attr.Sha, env.Attr.RunID,
		env.Attr.Groups,
		"--max-total-score", strconv.Itoa(conf.MaxTotalScore),
		skipIssueArg, skipScoreboardArg,
		skipFailedTableArg, submitterInIssueTitleArg,
	}
	stdoutBuf, err := runCommand(args)
	if err != nil {
		slog.Error("teapot run exec", "error", err)
		return
	}
	if json.Unmarshal(stdoutBuf.Bytes(), &runResult) != nil {
		slog.Error("unmarshal teapot result", "error", err,
			"stdout", stdoutBuf.String())
		return
	}
	slog.Info("teapot result", "result", runResult)
	return
}
