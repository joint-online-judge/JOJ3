package main

import (
	"encoding/json"
	"errors"
	"log/slog"
	"os"

	_ "focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/internal/executors"
	_ "focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/internal/parsers"
	"focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/internal/stage"

	"github.com/go-git/go-git/v5"
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

func generateStages(conf Conf) []stage.Stage {
	stages := []stage.Stage{}
	for _, s := range conf.Stages {
		var cmds []stage.Cmd
		defaultCmd := s.Executor.With.Default
		for _, optionalCmd := range s.Executor.With.Cases {
			cmd := s.Executor.With.Default
			err := copier.Copy(&cmd, &optionalCmd)
			if err != nil {
				slog.Error("generate stages", "error", err)
				os.Exit(1)
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
	return stages
}

func outputResult(conf Conf, results []stage.StageResult) error {
	content, err := json.Marshal(results)
	if err != nil {
		return err
	}
	return os.WriteFile(conf.OutputPath,
		append(content, []byte("\n")...), 0o600)
}

func main() {
	conf, err := commitMsgToConf()
	if err != nil {
		// FIXME: just for local testing purpose
		if errors.Is(err, git.ErrRepositoryNotExists) {
			conf = parseConfFile("conf.toml")
		} else {
			slog.Error("no conf found", "error", err)
			os.Exit(0)
		}
	}
	setupSlog(conf)
	stages := generateStages(conf)
	defer stage.Cleanup()
	results := stage.Run(stages)
	if err := outputResult(conf, results); err != nil {
		slog.Error("output result", "error", err)
	}
	if err := stage.Submit(
		conf.GiteaUrl,
		conf.GiteaToken,
		conf.GiteaOwner,
		conf.GiteaRepo,
		results,
	); err != nil {
		slog.Error("submit result", "error", err)
	}
}
