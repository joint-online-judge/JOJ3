package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/joint-online-judge/JOJ3/internal/executors"
	_ "github.com/joint-online-judge/JOJ3/internal/parsers"
	"github.com/joint-online-judge/JOJ3/internal/stage"

	"github.com/go-git/go-git/v5"
	"github.com/jinzhu/copier"
)

func getCommitMsg() (msg string, err error) {
	r, err := git.PlainOpen(".")
	if err != nil {
		return
	}
	ref, err := r.Head()
	if err != nil {
		return
	}
	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		return
	}
	msg = commit.Message
	return
}

func generateStages(conf Conf, group string) ([]stage.Stage, error) {
	stages := []stage.Stage{}
	for _, s := range conf.Stages {
		if s.Group != "" && group != "" && group != s.Group {
			continue
		}
		var cmds []stage.Cmd
		defaultCmd := s.Executor.With.Default
		for _, optionalCmd := range s.Executor.With.Cases {
			cmd := s.Executor.With.Default
			err := copier.Copy(&cmd, &optionalCmd)
			if err != nil {
				slog.Error("generate stages", "error", err)
				return stages, err
			}
			// since these 3 values are pointers, copier will always copy
			// them, so we need to check them manually
			if defaultCmd.Stdin != nil && optionalCmd.Stdin == nil {
				cmd.Stdin = defaultCmd.Stdin
			}
			if defaultCmd.Stdout != nil && optionalCmd.Stdout == nil {
				cmd.Stdout = defaultCmd.Stdout
			}
			if defaultCmd.Stderr != nil && optionalCmd.Stderr == nil {
				cmd.Stderr = defaultCmd.Stderr
			}
			cmds = append(cmds, cmd)
		}
		if len(s.Executor.With.Cases) == 0 {
			cmds = []stage.Cmd{defaultCmd}
		}
		stages = append(stages, stage.Stage{
			Name:         s.Name,
			ExecutorName: s.Executor.Name,
			ExecutorCmds: cmds,
			ParserName:   s.Parser.Name,
			ParserConf:   s.Parser.With,
		})
	}
	return stages, nil
}

func outputResult(outputPath string, results []stage.StageResult) error {
	content, err := json.Marshal(results)
	if err != nil {
		return err
	}
	return os.WriteFile(outputPath,
		append(content, []byte("\n")...), 0o600)
}

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
		msg, err = getCommitMsg()
		if err != nil {
			slog.Error("get commit msg", "error", err)
			return err
		}
	}
	conf, group, err := parseMsg(confRoot, confName, msg)
	if err != nil {
		slog.Error("parse msg", "error", err)
		return err
	}
	if err := setupSlog(conf.LogPath); err != nil { // after conf is loaded
		return err
	}
	slog.Info("debug log", "path", conf.LogPath)
	slog.Debug("conf loaded", "conf", conf)
	executors.InitWithConf(conf.SandboxExecServer, conf.SandboxToken)
	stages, err := generateStages(conf, group)
	if err != nil {
		slog.Error("generate stages", "error", err)
		return err
	}
	slog.Debug("stages generated", "stages", stages)
	defer stage.Cleanup()
	results, err := stage.Run(stages)
	if err != nil {
		slog.Error("run stages", "error", err)
		return err
	}
	slog.Debug("stages run done", "results", results)
	if err := outputResult(conf.OutputPath, results); err != nil {
		slog.Error("output result", "error", err)
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
