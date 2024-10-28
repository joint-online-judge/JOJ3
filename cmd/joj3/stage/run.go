package stage

import (
	"log/slog"

	"github.com/joint-online-judge/JOJ3/cmd/joj3/conf"
	executors "github.com/joint-online-judge/JOJ3/internal/executor"
	_ "github.com/joint-online-judge/JOJ3/internal/parser"
	"github.com/joint-online-judge/JOJ3/internal/stage"

	"github.com/jinzhu/copier"
)

func generateStages(conf *conf.Conf, group string) ([]stage.Stage, error) {
	stages := []stage.Stage{}
	existNames := map[string]bool{}
	for _, s := range conf.Stage.Stages {
		if s.Group != "" && group != s.Group {
			continue
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

func Run(conf *conf.Conf, group string) (
	stageResults []stage.StageResult, forceQuit bool, err error,
) {
	stageResultsOnError := []stage.StageResult{
		{
			Name: "Internal Error",
			Results: []stage.ParserResult{{
				Score:   0,
				Comment: "JOJ3 internal error, check the log in Gitea Actions.",
			}},
			ForceQuit: true,
		},
	}
	executors.InitWithConf(
		conf.Stage.SandboxExecServer,
		conf.Stage.SandboxToken,
	)
	stages, err := generateStages(conf, group)
	if err != nil {
		slog.Error("generate stages", "error", err)
		stageResults = stageResultsOnError
		forceQuit = true
		return
	}
	defer stage.Cleanup()
	stageResults, forceQuit, err = stage.Run(stages)
	if err != nil {
		slog.Error("run stages", "error", err)
		stageResults = stageResultsOnError
		forceQuit = true
		return
	}
	return
}
