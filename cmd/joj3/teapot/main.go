package teapot

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"

	"github.com/joint-online-judge/JOJ3/cmd/joj3/conf"
)

func Run(conf conf.Conf) error {
	if conf.SkipTeapot {
		return nil
	}
	os.Setenv("LOG_FILE_PATH", "/home/tt/.cache/joint-teapot-debug.log")
	os.Setenv("_TYPER_STANDARD_TRACEBACK", "1")
	envFilePath := "/home/tt/.config/teapot/teapot.env"
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
	cmd := exec.Command("joint-teapot", "joj3-scoreboard",
		envFilePath, conf.OutputPath, actor, conf.GradingRepoName, repoName,
		runNumber, conf.ScoreboardPath) // #nosec G204
	output, err := cmd.CombinedOutput()
	slog.Info("joint-teapot joj3-scoreboard", "output", string(output))
	if err != nil {
		slog.Error("joint-teapot joj3-scoreboard", "err", err)
		return err
	}
	cmd = exec.Command("joint-teapot", "joj3-failed-table",
		envFilePath, conf.OutputPath, actor, conf.GradingRepoName, repoName,
		runNumber, conf.FailedTablePath) // #nosec G204
	output, err = cmd.CombinedOutput()
	slog.Info("joint-teapot joj3-failed-table", "output", string(output))
	if err != nil {
		slog.Error("joint-teapot joj3-failed-table", "err", err)
		return err
	}
	cmd = exec.Command("joint-teapot", "joj3-create-result-issue",
		envFilePath, conf.OutputPath, repoName, runNumber) // #nosec G204
	output, err = cmd.CombinedOutput()
	slog.Info("joint-teapot joj3-create-result-issue", "output", string(output))
	if err != nil {
		slog.Error("joint-teapot joj3-create-result-issue", "err", err)
		return err
	}
	return nil
}
