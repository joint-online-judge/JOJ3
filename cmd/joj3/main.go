package main

import (
	"encoding/json"
	"flag"
	"log/slog"
	"os"

	"focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/internal/executors"
	_ "focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/internal/parsers"
	"focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/internal/stage"

	"github.com/go-git/go-git/v5"
	"github.com/jinzhu/copier"
)

func setupSlog(logLevel int) {
	lvl := new(slog.LevelVar)
	lvl.Set(slog.Level(logLevel))
	opts := &slog.HandlerOptions{Level: lvl}
	handler := slog.NewTextHandler(os.Stderr, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)
}

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

func generateStages(conf Conf) ([]stage.Stage, error) {
	stages := []stage.Stage{}
	for _, s := range conf.Stages {
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
		slog.Debug("parse stages conf", "cmds", cmds)
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
	metaConfPath string
	msg          string
)

func init() {
	flag.StringVar(&metaConfPath, "meta-conf", "meta-conf.json", "meta config file path")
	flag.StringVar(&msg, "msg", "", "message to trigger the running, leave empty to use git commit message on HEAD")
}

func mainImpl() error {
	setupSlog(int(slog.LevelInfo)) // before conf is loaded
	flag.Parse()
	if msg == "" {
		var err error
		msg, err = getCommitMsg()
		if err != nil {
			slog.Error("get commit msg", "error", err)
			return err
		}
	}
	conf, err := msgToConf(metaConfPath, msg)
	if err != nil {
		slog.Error("no conf found", "error", err)
		return err
	}
	setupSlog(conf.LogLevel) // after conf is loaded
	executors.InitWithConf(conf.SandboxExecServer, conf.SandboxToken)
	stages, err := generateStages(conf)
	if err != nil {
		slog.Error("generate stages", "error", err)
		return err
	}
	defer stage.Cleanup()
	results, err := stage.Run(stages)
	if err != nil {
		slog.Error("run stages", "error", err)
		return err
	}
	if err := outputResult(conf.OutputPath, results); err != nil {
		slog.Error("output result", "error", err)
		return err
	}
	return nil
}

func main() {
	if err := mainImpl(); err != nil {
		os.Exit(1)
	}
}
