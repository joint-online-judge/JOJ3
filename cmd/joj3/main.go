// Package main provides a joj3 executable, which runs various stages based on
// configuration files and commit message. The output is a JSON file.
package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"runtime/debug"
	"strings"

	joj3Conf "github.com/joint-online-judge/JOJ3/cmd/joj3/conf"
	"github.com/joint-online-judge/JOJ3/cmd/joj3/env"
	"github.com/joint-online-judge/JOJ3/internal/stage"
)

var (
	confFileRoot         string
	confFileName         string
	fallbackConfFileName string
	tag                  string
	printVersion         *bool
	Version              string = "debug"
)

func init() {
	flag.StringVar(&confFileRoot, "conf-root", ".", "root path for all config files")
	flag.StringVar(&confFileName, "conf-name", "conf.json", "filename for config files")
	flag.StringVar(&fallbackConfFileName, "fallback-conf-name", "", "filename for the fallback config file in conf-root, leave empty to use conf-name")
	flag.StringVar(&tag, "tag", "", "tag to trigger the running, when non-empty, should equal to the scope in msg")
	printVersion = flag.Bool("version", false, "print current version")
}

func getCommitMsg() (string, error) {
	commitMsg, err := joj3Conf.GetCommitMsg()
	if err != nil {
		slog.Error("get commit msg", "error", err)
		return "", err
	}
	env.Attr.CommitMsg = commitMsg
	return commitMsg, nil
}

func getConf(commitMsg string) (*joj3Conf.Conf, *joj3Conf.ConventionalCommit, error) {
	confPath, confStat, conventionalCommit, err := getConfPath(commitMsg)
	if err != nil {
		return nil, nil, err
	}
	conf, err := loadConf(confPath)
	if err != nil {
		return nil, nil, err
	}
	env.Attr.ConfName = conf.Name
	env.Attr.OutputPath = conf.Stage.OutputPath
	if err := showConfStat(confPath, confStat); err != nil {
		return nil, nil, err
	}
	return conf, conventionalCommit, nil
}

func getConfPath(commitMsg string) (string, fs.FileInfo, *joj3Conf.ConventionalCommit, error) {
	confPath, confStat, conventionalCommit, err := joj3Conf.GetConfPath(
		confFileRoot, confFileName, fallbackConfFileName, commitMsg, tag,
	)
	if err != nil {
		slog.Error("get conf path", "error", err)
		return "", nil, nil, err
	}
	slog.Info("try to load conf", "path", confPath)
	return confPath, confStat, conventionalCommit, nil
}

func loadConf(confPath string) (*joj3Conf.Conf, error) {
	conf, err := joj3Conf.ParseConfFile(confPath)
	if err != nil {
		slog.Error("parse conf", "error", err)
		return nil, err
	}
	slog.Debug("conf loaded", "conf", conf, "joj3 version", Version)
	return conf, nil
}

func showConfStat(confPath string, confStat fs.FileInfo) error {
	confSHA256, err := joj3Conf.GetSHA256(confPath)
	if err != nil {
		slog.Error("get sha256", "error", err)
		return err
	}
	slog.Info("conf info", "sha256", confSHA256, "modTime", confStat.ModTime(), "size", confStat.Size())
	return nil
}

func validateConf(conf *joj3Conf.Conf) error {
	if err := joj3Conf.CheckValid(conf); err != nil {
		slog.Error("conf not valid now", "error", err)
		return err
	}
	return nil
}

func run(conf *joj3Conf.Conf, conventionalCommit *joj3Conf.ConventionalCommit) error {
	groups := joj3Conf.MatchGroups(conf, conventionalCommit)
	env.Attr.Groups = strings.Join(groups, ",")
	env.Set()
	_, forceQuitStageName, err := runStages(
		conf,
		groups,
		func(stageResults []stage.StageResult, forceQuitStageName string) {
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

func mainImpl() (err error) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	slog.SetDefault(logger)
	flag.Parse()
	if *printVersion {
		fmt.Println(Version)
		return nil
	}
	if fallbackConfFileName == "" {
		fallbackConfFileName = confFileName
	}
	slog.Info("start joj3", "version", Version)
	commitMsg, err := getCommitMsg()
	if err != nil {
		return err
	}
	conf, conventionalCommit, err := getConf(commitMsg)
	if err != nil {
		return err
	}
	if err := setupSlog(conf); err != nil {
		return err
	}
	if err := validateConf(conf); err != nil {
		return err
	}
	if err := run(conf, conventionalCommit); err != nil {
		return err
	}
	return nil
}

func main() {
	var err error
	exitCode := 0
	defer func() {
		if r := recover(); r != nil {
			slog.Error(
				"panic recovered",
				"panic", r,
				"stack", string(debug.Stack()),
			)
			exitCode = 2
		}
		if err != nil {
			slog.Error("main exit", "error", err)
			exitCode = 1
		}
		if !*printVersion && exitCode == 0 {
			slog.Info("main exit", "status", "success")
		}
		os.Exit(exitCode)
	}()
	err = mainImpl()
}
