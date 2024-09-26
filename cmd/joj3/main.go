package main

import (
	"encoding/json"
	"flag"
	"log/slog"
	"os"

	"focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/internal/executors"
	_ "focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/internal/parsers"
	"focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/internal/stage"

	"github.com/jinzhu/copier"
)

func setupSlog(conf Conf) {
	lvl := new(slog.LevelVar)
	lvl.Set(slog.Level(conf.LogLevel))
	opts := &slog.HandlerOptions{Level: lvl}
	handler := slog.NewTextHandler(os.Stderr, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)
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

func outputResult(conf Conf, results []stage.StageResult) error {
	content, err := json.Marshal(results)
	if err != nil {
		return err
	}
	return os.WriteFile(conf.OutputPath,
		append(content, []byte("\n")...), 0o600)
}

var metaConfPath string

func init() {
	flag.StringVar(&metaConfPath, "meta-conf", "meta-conf.json", "meta config file path")
}

func mainImpl() error {
	flag.Parse()
	conf, err := commitMsgToConf(metaConfPath)
	if err != nil {
		slog.Error("no conf found", "error", err)
		return err
	}
	setupSlog(conf)
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
	if err := outputResult(conf, results); err != nil {
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
