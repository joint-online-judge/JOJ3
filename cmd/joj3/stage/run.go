package stage

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/joint-online-judge/JOJ3/cmd/joj3/conf"
	executors "github.com/joint-online-judge/JOJ3/internal/executor"
	_ "github.com/joint-online-judge/JOJ3/internal/parser"
	"github.com/joint-online-judge/JOJ3/internal/stage"

	"github.com/jinzhu/copier"
)

type StageResult stage.StageResult

func generateStages(conf *conf.Conf, groups []string) ([]stage.Stage, error) {
	stages := []stage.Stage{}
	existNames := map[string]bool{}
	for _, s := range conf.Stage.Stages {
		if s.Group != "" {
			var ok bool
			loweredStageGroup := strings.ToLower(s.Group)
			for _, group := range groups {
				if group == loweredStageGroup {
					ok = true
					break
				}
			}
			if !ok {
				continue
			}
		}
		_, ok := existNames[s.Name] // check for existence
		if ok {
			continue
		}
		existNames[s.Name] = true
		var cmds []stage.Cmd
		defaultCmd := s.Executor.With.Default
		for _, optionalCmd := range s.Executor.With.Cases {
			cmd := s.Executor.With.Default
			err := copier.CopyWithOption(
				&cmd,
				&optionalCmd,
				copier.Option{DeepCopy: true},
			)
			if err != nil {
				slog.Error("generate stages", "error", err)
				return stages, err
			}
			// since these 3 values are pointers, copier will always copy
			// them, so we need to check them manually
			if defaultCmd.Stdin != nil && optionalCmd.Stdin == nil {
				var stdin stage.CmdFile
				err := copier.CopyWithOption(
					&stdin,
					defaultCmd.Stdin,
					copier.Option{DeepCopy: true},
				)
				if err != nil {
					slog.Error("generate stages", "error", err)
					return stages, err
				}
				cmd.Stdin = &stdin
			}
			if defaultCmd.Stdout != nil && optionalCmd.Stdout == nil {
				var stdout stage.CmdFile
				err := copier.CopyWithOption(
					&stdout,
					defaultCmd.Stdout,
					copier.Option{DeepCopy: true},
				)
				if err != nil {
					slog.Error("generate stages", "error", err)
					return stages, err
				}
				cmd.Stdout = &stdout
			}
			if defaultCmd.Stderr != nil && optionalCmd.Stderr == nil {
				var stderr stage.CmdFile
				err := copier.CopyWithOption(
					&stderr,
					defaultCmd.Stderr,
					copier.Option{DeepCopy: true},
				)
				if err != nil {
					slog.Error("generate stages", "error", err)
					return stages, err
				}
				cmd.Stderr = &stderr
			}
			cmds = append(cmds, cmd)
		}
		if len(s.Executor.With.Cases) == 0 {
			cmds = []stage.Cmd{defaultCmd}
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

func newErrorStageResults(err error) []stage.StageResult {
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
	}
}

func Run(conf *conf.Conf, groups []string) (
	stageResults []stage.StageResult, forceQuit bool, err error,
) {
	executors.InitWithConf(
		conf.Stage.SandboxExecServer,
		conf.Stage.SandboxToken,
	)
	stages, err := generateStages(conf, groups)
	if err != nil {
		slog.Error("generate stages", "error", err)
		stageResults = newErrorStageResults(err)
		forceQuit = true
		return
	}
	defer stage.Cleanup()
	stageResults, forceQuit, err = stage.Run(stages)
	if err != nil {
		slog.Error("run stages", "error", err)
		stageResults = newErrorStageResults(err)
		forceQuit = true
		return
	}
	return
}
