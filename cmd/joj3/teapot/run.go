package teapot

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"

	"github.com/joint-online-judge/JOJ3/cmd/joj3/conf"
)

func Run(conf *conf.Conf) error {
	actions := os.Getenv("GITHUB_ACTIONS")
	if actions != "true" {
		slog.Info("teapot exit", "GITHUB_ACTIONS", actions)
		return nil
	}
	os.Setenv("LOG_FILE_PATH", conf.Teapot.LogPath)
	os.Setenv("_TYPER_STANDARD_TRACEBACK", "1")
	envFilePath := "/home/tt/.config/teapot/teapot.env"
	// TODO: pass sha to joint-teapot
	// sha := os.Getenv("GITHUB_SHA")
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
		outputBytes, err := cmd.CombinedOutput()
		output := re.ReplaceAllString(string(outputBytes), "")
		for _, line := range strings.Split(output, "\n") {
			if line == "" {
				continue
			}
			slog.Info(fmt.Sprintf("%s %s", name, cmdArgs[0]), "output", line)
		}
		return err
	}
	var wg sync.WaitGroup
	var scoreboardErr, failedTableErr, issueErr error
	wg.Add(2)
	go func() {
		defer wg.Done()
		if !conf.Teapot.SkipScoreboard {
			err := execCommand("joint-teapot", []string{
				"joj3-scoreboard", envFilePath, conf.Stage.OutputPath, actor,
				conf.Teapot.GradingRepoName, repoName, runNumber,
				conf.Teapot.ScoreboardPath, conf.Name,
			})
			if err != nil {
				scoreboardErr = err
			}
		}
		if !conf.Teapot.SkipFailedTable {
			err := execCommand("joint-teapot", []string{
				"joj3-failed-table", envFilePath, conf.Stage.OutputPath, actor,
				conf.Teapot.GradingRepoName, repoName, runNumber,
				conf.Teapot.FailedTablePath, conf.Name,
			})
			if err != nil {
				failedTableErr = err
			}
		}
	}()
	go func() {
		defer wg.Done()
		if !conf.Teapot.SkipIssue {
			err := execCommand("joint-teapot", []string{
				"joj3-create-result-issue", envFilePath, conf.Stage.OutputPath,
				repoName, runNumber, conf.Name,
			})
			if err != nil {
				issueErr = err
			}
		}
	}()
	wg.Wait()
	if scoreboardErr != nil || failedTableErr != nil || issueErr != nil {
		slog.Error("teapot exit", "scoreboardErr", scoreboardErr,
			"failedTableErr", failedTableErr, "issueErr", issueErr)
		return fmt.Errorf("teapot exit")
	}
	return nil
}
