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
	confRoot    string
	confName    string
	msg         string
	showVersion *bool
	Version     string = "debug"
)

func init() {
	flag.StringVar(&confRoot, "conf-root", ".", "root path for all config files")
	flag.StringVar(&confName, "conf-name", "conf.json", "filename for config files")
	flag.StringVar(&msg, "msg", "", "message to trigger the running, leave empty to use git commit message on HEAD")
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
	slog.Info("start joj3", "version", Version)
	if msg == "" {
		var err error
		msg, err = conf.GetCommitMsg()
		if err != nil {
			slog.Error("get commit msg", "error", err)
			return err
		}
	}
	confPath, group, err := conf.ParseMsg(confRoot, confName, msg)
	if err != nil {
		slog.Error("parse msg", "error", err)
		validScopes, scopeErr := conf.ListValidScopes(
			confRoot, confName)
		if scopeErr != nil {
			slog.Error("list valid scopes", "error", scopeErr)
			return err
		}
		slog.Info("HINT: use valid scopes in commit message",
			"valid scopes", validScopes)
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
	if err := stage.Run(confObj, group); err != nil {
		slog.Error("stage run", "error", err)
		return err
	}
	if err := teapot.Run(confObj); err != nil {
		slog.Error("teapot run", "error", err)
		return err
	}
	return nil
}

func main() {
	if err := mainImpl(); err != nil {
		slog.Error("main exit", "error", err)
		os.Exit(1)
	}
}
