package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"

	"github.com/joint-online-judge/JOJ3/cmd/joj3/env"
	"github.com/joint-online-judge/JOJ3/internal/conf"
)

type CheckResult struct {
	Name        string `json:"name"`
	SubmitCount int    `json:"submit_count"`
	MaxCount    int    `json:"max_count"`
	TimePeriod  int    `json:"time_period"`
}

func check(conf *conf.Conf) (checkResults []CheckResult, err error) {
	os.Setenv("LOG_FILE_PATH", conf.Teapot.LogPath)
	os.Setenv("_TYPER_STANDARD_TRACEBACK", "1")
	actor := os.Getenv("GITHUB_ACTOR")
	repository := os.Getenv("GITHUB_REPOSITORY")
	if actor == "" ||
		repository == "" ||
		strings.Count(repository, "/") != 1 {
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
		actor, conf.Teapot.GradingRepoName, repoName,
		conf.Teapot.ScoreboardPath, conf.Name,
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
	if json.Unmarshal(stdoutBuf.Bytes(), &checkResults) != nil {
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
			"in last %d hour(s): submit count %d, max count %d\n",
			checkResult.TimePeriod,
			checkResult.SubmitCount,
			checkResult.MaxCount,
		)
		if useGroup && checkResult.SubmitCount+1 > checkResult.MaxCount {
			err = fmt.Errorf("submit count exceeded")
		}
	}
	return
}

func setupSlog() {
	opts := &slog.HandlerOptions{}
	handler := slog.NewTextHandler(os.Stderr, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)
}

var (
	confPath string
	Version  string = "debug"
)

func mainImpl() (err error) {
	showVersion := flag.Bool("version", false, "print current version")
	flag.StringVar(&confPath, "conf-path", "./conf.json", "path for config file")
	flag.Parse()
	if *showVersion {
		fmt.Println(Version)
		return
	}
	setupSlog()
	slog.Info("start teapot-checker", "version", Version)
	commitMsg, err := conf.GetCommitMsg()
	if err != nil {
		slog.Error("get commit msg", "error", err)
		return
	}
	conventionalCommit, err := conf.ParseConventionalCommit(commitMsg)
	if err != nil {
		slog.Error("parse commit msg", "error", err)
		return
	}
	confObj, _, err := conf.ParseConfFile(confPath)
	if err != nil {
		slog.Error("parse conf", "error", err)
		return
	}
	groups := conf.MatchGroups(confObj, conventionalCommit)
	checkResults, err := check(confObj)
	if err != nil {
		slog.Error("teapot check", "error", err)
		return
	}
	output, err := generateOutput(checkResults, groups)
	if err != nil {
		slog.Error("generate output", "error", err)
		return
	}
	fmt.Println(output)
	return
}

func main() {
	if err := mainImpl(); err != nil {
		slog.Error("main exit", "error", err)
		os.Exit(1)
	}
}
