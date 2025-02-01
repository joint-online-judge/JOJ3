package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/joint-online-judge/JOJ3/cmd/joj3/env"
	"github.com/joint-online-judge/JOJ3/cmd/joj3/stage"
	"github.com/joint-online-judge/JOJ3/internal/conf"
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

	if err := setupSlog(confObj); err != nil { // before conf is loaded
		slog.Error("setup slog", "error", err)
		return err
	}

	// parse flag & conf file
	flag.Parse()
	if *showVersion {
		fmt.Println(Version)
		return nil
	}
	if fallbackConfFileName == "" {
		fallbackConfFileName = confFileName
	}
	slog.Info("start joj3", "version", Version)
	commitMsg, err := conf.GetCommitMsg()
	if err != nil {
		slog.Error("get commit msg", "error", err)
		return err
	}
	env.Attr.CommitMsg = commitMsg
	confPath, confStat, conventionalCommit, err := conf.GetConfPath(
		confFileRoot, confFileName, fallbackConfFileName, commitMsg, tag)
	if err != nil {
		slog.Error("get conf path", "error", err)
		return err
	}
	slog.Info("try to load conf", "path", confPath)
	confObj, err = conf.ParseConfFile(confPath)
	if err != nil {
		slog.Error("parse conf", "error", err)
		return err
	}
	env.Attr.ConfName = confObj.Name
	slog.Debug("conf loaded", "conf", confObj)
	if err := setupSlog(confObj); err != nil { // after conf is loaded
		slog.Error("setup slog", "error", err)
		return err
	}

	// log conf file info
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

	// run stages
	groups := conf.MatchGroups(confObj, conventionalCommit)
	env.Attr.Groups = strings.Join(groups, ",")
	_, forceQuitStageName, err := stage.Run(
		confObj,
		groups,
		func(
			stageResults []internalStage.StageResult,
			forceQuitStageName string,
		) {
			env.Attr.ForceQuitStageName = forceQuitStageName
			env.Set()
		},
	)
	if err != nil {
		slog.Error("stage run", "error", err)
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
