package teapot

import (
	"bufio"
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
		if err := cmd.Start(); err != nil {
			slog.Error("cmd start", "error", err)
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
			return err
		}
		wg.Wait()
		return err
	}
	var wg sync.WaitGroup
	var scoreboardErr, failedTableErr, issueErr error
	if !conf.Teapot.SkipScoreboard {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := execCommand("joint-teapot", []string{
				"joj3-scoreboard", envFilePath, conf.Stage.OutputPath, actor,
				conf.Teapot.GradingRepoName, repoName, runNumber,
				conf.Teapot.ScoreboardPath, conf.Name, sha,
			})
			if err != nil {
				scoreboardErr = err
			}
		}()
	}
	if !conf.Teapot.SkipFailedTable {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := execCommand("joint-teapot", []string{
				"joj3-failed-table", envFilePath, conf.Stage.OutputPath, actor,
				conf.Teapot.GradingRepoName, repoName, runNumber,
				conf.Teapot.FailedTablePath, conf.Name, sha,
			})
			if err != nil {
				failedTableErr = err
			}
		}()
	}
	if !conf.Teapot.SkipIssue {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := execCommand("joint-teapot", []string{
				"joj3-create-result-issue", envFilePath, conf.Stage.OutputPath,
				repoName, runNumber, conf.Name, actor, sha,
			})
			if err != nil {
				issueErr = err
			}
		}()
	}
	wg.Wait()
	if scoreboardErr != nil || failedTableErr != nil || issueErr != nil {
		slog.Error("teapot exit", "scoreboardErr", scoreboardErr,
			"failedTableErr", failedTableErr, "issueErr", issueErr)
		return fmt.Errorf("teapot exit")
	}
	return nil
}
