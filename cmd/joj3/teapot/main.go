package teapot

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"regexp"
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
	re := regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
	cmd := exec.Command("joint-teapot", "joj3-scoreboard",
		envFilePath, conf.OutputPath, actor, conf.GradingRepoName, repoName,
		runNumber, conf.ScoreboardPath, conf.Name) // #nosec G204
	outputBytes, err := cmd.CombinedOutput()
	output := re.ReplaceAllString(string(outputBytes), "")
	for _, line := range strings.Split(output, "\n") {
		if line == "" {
			continue
		}
		slog.Info("joint-teapot joj3-scoreboard", "output", line)
	}
	if err != nil {
		slog.Error("joint-teapot joj3-scoreboard", "err", err)
		return err
	}
	cmd = exec.Command("joint-teapot", "joj3-failed-table",
		envFilePath, conf.OutputPath, actor, conf.GradingRepoName, repoName,
		runNumber, conf.FailedTablePath, conf.Name) // #nosec G204
	outputBytes, err = cmd.CombinedOutput()
	output = re.ReplaceAllString(string(outputBytes), "")
	for _, line := range strings.Split(output, "\n") {
		if line == "" {
			continue
		}
		slog.Info("joint-teapot joj3-failed-table", "output", line)
	}
	if err != nil {
		slog.Error("joint-teapot joj3-failed-table", "err", err)
		return err
	}
	cmd = exec.Command("joint-teapot", "joj3-create-result-issue",
		envFilePath, conf.OutputPath, repoName, runNumber, conf.Name) // #nosec G204
	outputBytes, err = cmd.CombinedOutput()
	output = re.ReplaceAllString(string(outputBytes), "")
	for _, line := range strings.Split(output, "\n") {
		if line == "" {
			continue
		}
		slog.Info("joint-teapot joj3-create-result-issue", "output", line)
	}
	if err != nil {
		slog.Error("joint-teapot joj3-create-result-issue", "err", err)
		return err
	}
	return nil
}
