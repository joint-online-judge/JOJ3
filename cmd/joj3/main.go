package main

import (
	"flag"
	"fmt"
	"log/slog"

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

func main() {
	if err := setupSlog(""); err != nil { // before conf is loaded
		slog.Error("setup slog", "error", err)
		return
	}
	flag.Parse()
	if *showVersion {
		fmt.Println(Version)
		return
	}
	slog.Info("start joj3", "version", Version)
	if msg == "" {
		var err error
		msg, err = conf.GetCommitMsg()
		if err != nil {
			slog.Error("get commit msg", "error", err)
			return
		}
	}
	confObj, group, err := conf.ParseMsg(confRoot, confName, msg)
	if err != nil {
		slog.Error("parse msg", "error", err)
		validScopes, scopeErr := conf.ListValidScopes(
			confRoot, confName, msg)
		if scopeErr != nil {
			slog.Error("list valid scopes", "error", scopeErr)
			return
		}
		slog.Info("hint: valid scopes in commit message", "scopes", validScopes)
		return
	}
	if err := setupSlog(confObj.LogPath); err != nil { // after conf is loaded
		slog.Error("setup slog", "error", err)
		return
	}
	if err := stage.Run(confObj, group); err != nil {
		slog.Error("stage run", "error", err)
		return
	}
	if err := teapot.Run(confObj); err != nil {
		slog.Error("teapot run", "error", err)
		return
	}
}
