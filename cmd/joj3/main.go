package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/joint-online-judge/JOJ3/cmd/joj3/conf"
	"github.com/joint-online-judge/JOJ3/cmd/joj3/stage"
	"github.com/joint-online-judge/JOJ3/cmd/joj3/teapot"
)

var (
	confRoot         string
	confName         string
	fallbackConfName string
	tag              string
	msg              string
	showVersion      *bool
	Version          string = "debug"
)

func init() {
	flag.StringVar(&confRoot, "conf-root", ".", "root path for all config files")
	flag.StringVar(&confName, "conf-name", "conf.json", "filename for config files")
	flag.StringVar(&fallbackConfName, "fallback-conf-name", "", "filename for the fallback config file in conf-root, leave empty to use conf-name")
	flag.StringVar(&tag, "tag", "", "tag to trigger the running, when non-empty, should equal to the scope in msg")
	// TODO: remove this flag
	flag.StringVar(&msg, "msg", "", "[DEPRECATED] will be ignored")
	showVersion = flag.Bool("version", false, "print current version")
}

func mainImpl() error {
	if err := setupSlog(""); err != nil { // before conf is loaded
		slog.Error("setup slog", "error", err)
		return err
	}
	flag.Parse()
	if *showVersion {
		fmt.Println(Version)
		return nil
	}
	if fallbackConfName == "" {
		fallbackConfName = confName
	}
	slog.Info("start joj3", "version", Version)
	msg, err := conf.GetCommitMsg()
	if err != nil {
		slog.Error("get commit msg", "error", err)
		return err
	}
	confPath, confStat, conventionalCommit, err := conf.GetConfPath(
		confRoot, confName, fallbackConfName, msg, tag)
	if err != nil {
		slog.Error("get conf path", "error", err)
		return err
	}
	slog.Info("try to load conf", "path", confPath)
	confObj, err := conf.ParseConfFile(confPath)
	if err != nil {
		slog.Error("parse conf", "error", err)
		return err
	}
	slog.Debug("conf loaded", "conf", confObj)
	if err := setupSlog(confObj.LogPath); err != nil { // after conf is loaded
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
	stageResults, stageForceQuit, err := stage.Run(confObj, groups)
	if err != nil {
		slog.Error("stage run", "error", err)
	}
	if err = stage.Write(confObj.Stage.OutputPath, stageResults); err != nil {
		slog.Error("stage write", "error", err)
		return err
	}
	if err := teapot.Run(confObj); err != nil {
		slog.Error("teapot run", "error", err)
		return err
	}
	if stageForceQuit {
		slog.Info("stage force quit")
		return fmt.Errorf("stage force quit")
	}
	return nil
}

func main() {
	if err := mainImpl(); err != nil {
		slog.Error("main exit", "error", err)
		os.Exit(1)
	}
}
