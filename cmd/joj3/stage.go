package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/joint-online-judge/JOJ3/cmd/joj3/conf"
	"github.com/joint-online-judge/JOJ3/internal/executor"
	_ "github.com/joint-online-judge/JOJ3/internal/parser"
	"github.com/joint-online-judge/JOJ3/internal/stage"

	"github.com/jinzhu/copier"
)

type StageResult stage.StageResult

func newStageCmd(defaultCmd stage.Cmd, optionalCmd conf.OptionalCmd) (stage.Cmd, error) {
	var cmd stage.Cmd
	err := copier.CopyWithOption(
		&cmd,
		&defaultCmd,
		copier.Option{DeepCopy: true},
	)
	if err != nil {
		return cmd, err
	}
	err = copier.CopyWithOption(
		&cmd,
		&optionalCmd,
		copier.Option{DeepCopy: true, IgnoreEmpty: true},
	)
	if err != nil {
		return cmd, err
	}
	return cmd, nil
}

func generateStages(confStages []conf.ConfStage, groups []string) (
	[]stage.Stage, error,
) {
	stages := []stage.Stage{}
	existNames := map[string]bool{}
	for i, s := range confStages {
		if s.Name == "" {
			s.Name = fmt.Sprintf("stage-%d", i)
		}
		var ok bool
		if s.Group != "" {
			for _, group := range groups {
				if strings.EqualFold(group, s.Group) {
					ok = true
					break
				}
			}
		}
		if !ok && len(s.Groups) > 0 {
			for _, group := range groups {
				for _, g := range s.Groups {
					if strings.EqualFold(group, g) {
						ok = true
						break
					}
				}
				if ok {
					break
				}
			}
		}
		if !ok {
			continue
		}
		_, ok = existNames[s.Name] // check for existence
		if ok {
			continue
		}
		existNames[s.Name] = true
		var cmds []stage.Cmd
		defaultCmd := s.Executor.With.Default
		if len(s.Executor.With.Cases) == 0 {
			cmds = []stage.Cmd{defaultCmd}
		} else {
			for _, optionalCmd := range s.Executor.With.Cases {
				cmd, err := newStageCmd(defaultCmd, optionalCmd)
				if err != nil {
					slog.Error("generate stages", "error", err)
					return stages, err
				}
				cmds = append(cmds, cmd)
			}
		}
		parsers := []stage.StageParser{}
		for _, p := range s.Parsers {
			parsers = append(parsers, stage.StageParser{
				Name: p.Name,
				Conf: p.With,
			})
		}
		stages = append(stages, stage.Stage{
			Name: s.Name,
			Executor: stage.StageExecutor{
				Name: s.Executor.Name,
				Cmds: cmds,
			},
			Parsers: parsers,
		})
	}
	slog.Debug("stages generated", "stages", stages)
	return stages, nil
}

func newErrorStageResults(err error) ([]stage.StageResult, string) {
	return []stage.StageResult{
		{
			Name: "Internal Error",
			Results: []stage.ParserResult{{
				Score: 0,
				Comment: "JOJ3 internal error, " +
					"check the log in Gitea Actions.\n" +
					fmt.Sprintf("Error: `%s`", err),
			}},
			ForceQuit: true,
		},
	}, "Internal Error"
}

func runStages(
	conf *conf.Conf,
	groups []string,
	onStagesComplete func([]stage.StageResult, string),
) (
	stageResults []stage.StageResult, forceQuitStageName string, err error,
) {
	executor.InitWithConf(
		conf.Stage.SandboxExecServer,
		conf.Stage.SandboxToken,
	)
	preStages, err := generateStages(conf.Stage.PreStages, groups)
	if err != nil {
		slog.Error("generate preStages", "error", err)
		stageResults, forceQuitStageName = newErrorStageResults(err)
		return stageResults, forceQuitStageName, err
	}
	stages, err := generateStages(conf.Stage.Stages, groups)
	if err != nil {
		slog.Error("generate stages", "error", err)
		stageResults, forceQuitStageName = newErrorStageResults(err)
		return stageResults, forceQuitStageName, err
	}
	postStages, err := generateStages(conf.Stage.PostStages, groups)
	if err != nil {
		slog.Error("generate postStages", "error", err)
		stageResults, forceQuitStageName = newErrorStageResults(err)
		return stageResults, forceQuitStageName, err
	}
	defer stage.Cleanup()
	// ignore force quit in preStages & postStages
	slog.Info("run preStages")
	_, _, err = stage.Run(preStages)
	if err != nil {
		slog.Error("run preStages", "error", err)
	}
	slog.Info("run stages")
	stageResults, forceQuitStageName, err = stage.Run(stages)
	if err != nil {
		slog.Error("run stages", "error", err)
		stageResults, forceQuitStageName = newErrorStageResults(err)
	}
	onStagesComplete(stageResults, forceQuitStageName)
	slog.Info("output result start", "path", conf.Stage.OutputPath)
	slog.Debug("output result start",
		"path", conf.Stage.OutputPath, "results", stageResults)
	content, err := json.Marshal(stageResults)
	if err != nil {
		slog.Error("marshal stageResults", "error", err)
	}
	err = os.WriteFile(conf.Stage.OutputPath,
		append(content, []byte("\n")...), 0o600)
	if err != nil {
		slog.Error("write stageResults", "error", err)
	}
	slog.Info("run postStages")
	_, _, err = stage.Run(postStages)
	if err != nil {
		slog.Error("run postStages", "error", err)
	}
	return stageResults, forceQuitStageName, err
}
