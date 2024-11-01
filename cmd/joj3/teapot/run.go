package teapot

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/joint-online-judge/JOJ3/cmd/joj3/conf"
)

func Run(conf *conf.Conf, runID string) error {
	os.Setenv("LOG_FILE_PATH", conf.Teapot.LogPath)
	os.Setenv("_TYPER_STANDARD_TRACEBACK", "1")
	sha := os.Getenv("GITHUB_SHA")
	actor := os.Getenv("GITHUB_ACTOR")
	repository := os.Getenv("GITHUB_REPOSITORY")
	runNumber := os.Getenv("GITHUB_RUN_NUMBER")
	if actor == "" || repository == "" || strings.Count(repository, "/") != 1 ||
		runNumber == "" {
		slog.Error("teapot env not set")
		return fmt.Errorf("teapot env not set")
	}
	repoParts := strings.Split(repository, "/")
	repoName := repoParts[1]
	re := regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
	execCommand := func(name string, cmdArgs []string) error {
		cmd := exec.Command(name, cmdArgs...) // #nosec G204
		stderr, err := cmd.StderrPipe()
		if err != nil {
			slog.Error("stderr pipe", "error", err)
			return err
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
				slog.Info(fmt.Sprintf("%s %s", name, cmdArgs[0]), "stderr", text)
			}
			wg.Done()
			if scanner.Err() != nil {
				slog.Error("stderr scanner", "error", scanner.Err())
			}
		}()
		if err = cmd.Start(); err != nil {
			slog.Error("cmd start", "error", err)
			return err
		}
		wg.Wait()
		return err
	}
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
	if err := execCommand("joint-teapot", []string{
		"joj3-all", conf.Teapot.EnvFilePath, conf.Stage.OutputPath, actor,
		conf.Teapot.GradingRepoName, repoName, runNumber,
		conf.Teapot.ScoreboardPath, conf.Teapot.FailedTablePath,
		conf.Name, sha, runID, strconv.Itoa(conf.Teapot.MaxTotalScore),
		skipIssueArg, skipScoreboardArg,
		skipFailedTableArg, submitterInIssueTitleArg,
	}); err != nil {
		slog.Error("teapot exit", "error", err)
		return fmt.Errorf("teapot exit")
	}
	return nil
}
