package teapot

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/joint-online-judge/JOJ3/cmd/joj3/conf"
	"github.com/joint-online-judge/JOJ3/cmd/joj3/env"
)

type TeapotResult struct {
	Issue  int    `json:"issue"`
	Action int    `json:"action"`
	Sha    string `json:"sha"`
}

func Run(conf *conf.Conf, groups []string) (
	teapotResult TeapotResult, err error,
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
	re := regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
	cmd := exec.Command("joint-teapot",
		"joj3-all", conf.Teapot.EnvFilePath, conf.Stage.OutputPath,
		env.Attr.Actor, conf.Teapot.GradingRepoName, repoName,
		env.Attr.RunNumber, conf.Teapot.ScoreboardPath,
		conf.Teapot.FailedTablePath,
		conf.Name, env.Attr.Sha, env.Attr.RunID,
		strings.Join(groups, ","),
		"--max-total-score", strconv.Itoa(conf.MaxTotalScore),
		skipIssueArg, skipScoreboardArg,
		skipFailedTableArg, submitterInIssueTitleArg,
	) // #nosec G204
	stdoutBuf := new(bytes.Buffer)
	cmd.Stdout = stdoutBuf
	stderr, err := cmd.StderrPipe()
	if err != nil {
		slog.Error("stderr pipe", "error", err)
		return
	}
	var wg sync.WaitGroup
	wg.Add(1)
	scanner := bufio.NewScanner(stderr)
	go func() {
		for scanner.Scan() {
			text := re.ReplaceAllString(scanner.Text(), "")
			if text == "" {
				continue
			}
			slog.Info("joint-teapot", "stderr", text)
		}
		wg.Done()
		if scanner.Err() != nil {
			slog.Error("stderr scanner", "error", scanner.Err())
		}
	}()
	if err = cmd.Start(); err != nil {
		slog.Error("cmd start", "error", err)
		return
	}
	wg.Wait()
	if err = cmd.Wait(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode := exitErr.ExitCode()
			slog.Error("cmd completed with non-zero exit code",
				"error", err,
				"exitCode", exitCode)
		} else {
			slog.Error("cmd wait", "error", err)
		}
		return
	}
	if json.Unmarshal(stdoutBuf.Bytes(), &teapotResult) != nil {
		slog.Error("unmarshal teapot result", "error", err,
			"stdout", stdoutBuf.String())
		return
	}
	slog.Info("teapot result", "result", teapotResult)
	return
}
