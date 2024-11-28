package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/joint-online-judge/JOJ3/cmd/joj3/conf"
	"github.com/joint-online-judge/JOJ3/cmd/joj3/env"
	"github.com/joint-online-judge/JOJ3/cmd/joj3/stage"
	"github.com/joint-online-judge/JOJ3/cmd/joj3/teapot"
	internalStage "github.com/joint-online-judge/JOJ3/internal/stage"
)

var (
	confFileRoot         string
	confFileName         string
	fallbackConfFileName string
	tag                  string
	showVersion          *bool
	Version              string = "debug"
)

func init() {
	flag.StringVar(&confFileRoot, "conf-root", ".", "root path for all config files")
	flag.StringVar(&confFileName, "conf-name", "conf.json", "filename for config files")
	flag.StringVar(&fallbackConfFileName, "fallback-conf-name", "", "filename for the fallback config file in conf-root, leave empty to use conf-name")
	flag.StringVar(&tag, "tag", "", "tag to trigger the running, when non-empty, should equal to the scope in msg")
	showVersion = flag.Bool("version", false, "print current version")
}

func mainImpl() (err error) {
	confObj := new(conf.Conf)
	var stageResults []internalStage.StageResult
	var forceQuitStageName string
	var teapotCheckResults []teapot.CheckResult
	var teapotRunResult teapot.RunResult
	var commitMsg string
	defer func() {
		totalScore := 0
		for _, stageResult := range stageResults {
			for _, result := range stageResult.Results {
				totalScore += result.Score
			}
		}
		cappedTotalScore := totalScore
		if confObj.MaxTotalScore >= 0 {
			cappedTotalScore = min(totalScore, confObj.MaxTotalScore)
		}
		slog.Info(
			"joj3 summary",
			"totalScore", totalScore,
			"cappedTotalScore", cappedTotalScore,
			"forceQuit", forceQuitStageName != "",
			"forceQuitStageName", forceQuitStageName,
			"issue", teapotRunResult.Issue,
			"action", teapotRunResult.Action,
			"sha", teapotRunResult.Sha,
			"commitMsg", commitMsg,
			"error", err,
		)
	}()
	if err := setupSlog(confObj); err != nil { // before conf is loaded
		slog.Error("setup slog", "error", err)
		return err
	}
	flag.Parse()
	if *showVersion {
		fmt.Println(Version)
		return nil
	}
	if fallbackConfFileName == "" {
		fallbackConfFileName = confFileName
	}
	slog.Info("start joj3", "version", Version)
	commitMsg, err = conf.GetCommitMsg()
	if err != nil {
		slog.Error("get commit msg", "error", err)
		return err
	}
	confPath, confStat, conventionalCommit, err := conf.GetConfPath(
		confFileRoot, confFileName, fallbackConfFileName, commitMsg, tag)
	if err != nil {
		slog.Error("get conf path", "error", err)
		return err
	}
	slog.Info("try to load conf", "path", confPath)
	confObj, confName, err := conf.ParseConfFile(confPath)
	if err != nil {
		slog.Error("parse conf", "error", err)
		return err
	}
	env.Attr.ConfName = confName
	slog.Debug("conf loaded", "conf", confObj)
	if err := setupSlog(confObj); err != nil { // after conf is loaded
		slog.Error("setup slog", "error", err)
		return err
	}
	sha256, err := conf.GetSHA256(confPath)
	if err != nil {
		slog.Error("get sha256", "error", err)
		return err
	}
	slog.Info("conf info", "sha256", sha256, "modTime", confStat.ModTime(),
		"size", confStat.Size())
	if err := conf.CheckExpire(confObj); err != nil {
		slog.Error("conf check expire", "error", err)
		return err
	}
	groups := conf.MatchGroups(confObj, conventionalCommit)
	if len(confObj.Teapot.Groups) != 0 {
		teapotCheckResults, err = teapot.Check(confObj)
		if err != nil {
			slog.Error("teapot check", "error", err)
		}
	} else {
		slog.Info("teapot check disabled")
	}
	stageResults, forceQuitStageName, err = stage.Run(
		confObj, groups, teapotCheckResults,
	)
	if err != nil {
		slog.Error("stage run", "error", err)
	}
	if err = stage.Write(confObj.Stage.OutputPath, stageResults); err != nil {
		slog.Error("stage write", "error", err)
		return err
	}
	teapotRunResult, err = teapot.Run(confObj, groups)
	if err != nil {
		slog.Error("teapot run", "error", err)
		return err
	}
	if forceQuitStageName != "" {
		slog.Info("stage force quit", "name", forceQuitStageName)
		return fmt.Errorf("stage force quit with name %s", forceQuitStageName)
	}
	return nil
}

func main() {
	if err := mainImpl(); err != nil {
		slog.Error("main exit", "error", err)
		os.Exit(1)
	}
}
