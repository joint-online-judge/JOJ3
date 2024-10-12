package conf

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/joint-online-judge/JOJ3/internal/stage"
	"github.com/koding/multiconfig"
)

type ConfStage struct {
	Name     string
	Group    string
	Executor struct {
		Name string
		With struct {
			Default stage.Cmd
			Cases   []OptionalCmd
		}
	}
	Parsers []struct {
		Name string
		With interface{}
	}
}

type Conf struct {
	Name    string `default:"unknown"`
	LogPath string `default:""`
	Stage   struct {
		SandboxExecServer string `default:"localhost:5051"`
		SandboxToken      string `default:""`
		OutputPath        string `default:"joj3_result.json"`
		Stages            []ConfStage
	}
	Teapot struct {
		LogPath         string `default:"/home/tt/.cache/joint-teapot-debug.log"`
		ScoreboardPath  string `default:"scoreboard.csv"`
		FailedTablePath string `default:"failed-table.md"`
		GradingRepoName string `default:""`
		SkipIssue       bool   `default:"false"`
		SkipScoreboard  bool   `default:"false"`
		SkipFailedTable bool   `default:"false"`
	}
	// TODO: remove the following backward compatibility fields
	SandboxExecServer string `default:"localhost:5051"`
	SandboxToken      string `default:""`
	OutputPath        string `default:"joj3_result.json"`
	GradingRepoName   string `default:""`
	SkipTeapot        bool   `default:"true"`
	ScoreboardPath    string `default:"scoreboard.csv"`
	FailedTablePath   string `default:"failed-table.md"`
	Stages            []struct {
		Name     string
		Group    string
		Executor struct {
			Name string
			With struct {
				Default stage.Cmd
				Cases   []OptionalCmd
			}
		}
		Parser struct {
			Name string
			With interface{}
		}
	}
}

type OptionalCmd struct {
	Args   *[]string
	Env    *[]string
	Stdin  *stage.CmdFile
	Stdout *stage.CmdFile
	Stderr *stage.CmdFile

	CPULimit     *uint64
	RealCPULimit *uint64
	ClockLimit   *uint64
	MemoryLimit  *uint64
	StackLimit   *uint64
	ProcLimit    *uint64
	CPURateLimit *uint64
	CPUSetLimit  *string

	CopyIn       *map[string]stage.CmdFile
	CopyInCached *map[string]string
	CopyInDir    *string

	CopyOut       *[]string
	CopyOutCached *[]string
	CopyOutMax    *uint64
	CopyOutDir    *string

	TTY               *bool
	StrictMemoryLimit *bool
	DataSegmentLimit  *bool
	AddressSpaceLimit *bool
}

type ConventionalCommit struct {
	Type        string
	Scope       string
	Description string
	Body        string
	Footer      string
}

func GetCommitMsg() (msg string, err error) {
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

func parseConventionalCommit(commit string) (*ConventionalCommit, error) {
	re := regexp.MustCompile(`(?s)^(\w+)(\(([^)]+)\))?!?: (.+?)(\n\n(.+?))?(\n\n(.+))?$`)
	matches := re.FindStringSubmatch(strings.TrimSpace(commit))
	if matches == nil {
		return nil, fmt.Errorf("invalid conventional commit format")
	}
	cc := &ConventionalCommit{
		Type:        matches[1],
		Scope:       matches[3],
		Description: strings.TrimSpace(matches[4]),
		Body:        strings.TrimSpace(matches[6]),
		Footer:      strings.TrimSpace(matches[8]),
	}
	return cc, nil
}

func ParseConfFile(path string) (conf Conf, err error) {
	d := &multiconfig.DefaultLoader{}
	d.Loader = multiconfig.MultiLoader(
		&multiconfig.TagLoader{},
		&multiconfig.JSONLoader{Path: path},
	)
	d.Validator = multiconfig.MultiValidator(&multiconfig.RequiredValidator{})
	if err = d.Load(&conf); err != nil {
		slog.Error("parse stages conf", "error", err)
		return
	}
	if err = d.Validate(&conf); err != nil {
		slog.Error("validate stages conf", "error", err)
		return
	}
	// TODO: remove the following backward compatibility codes
	if len(conf.Stage.Stages) == 0 {
		conf.Stage.SandboxExecServer = conf.SandboxExecServer
		conf.Stage.SandboxToken = conf.SandboxToken
		conf.Stage.OutputPath = conf.OutputPath
		conf.Stage.Stages = make([]ConfStage, len(conf.Stages))
		for i, stage := range conf.Stages {
			conf.Stage.Stages[i].Name = stage.Name
			conf.Stage.Stages[i].Group = stage.Group
			conf.Stage.Stages[i].Executor = stage.Executor
			conf.Stage.Stages[i].Parsers = []struct {
				Name string
				With interface{}
			}{
				{
					Name: stage.Parser.Name,
					With: stage.Parser.With,
				},
			}
		}
		conf.Teapot.GradingRepoName = conf.GradingRepoName
		conf.Teapot.ScoreboardPath = conf.ScoreboardPath
		conf.Teapot.FailedTablePath = conf.FailedTablePath
		if conf.SkipTeapot {
			conf.Teapot.SkipScoreboard = true
			conf.Teapot.SkipFailedTable = true
			conf.Teapot.SkipIssue = true
		}
	}
	return
}

func ParseMsg(confRoot, confName, msg string) (confPath, group string, err error) {
	slog.Info("parse msg", "msg", msg)
	conventionalCommit, err := parseConventionalCommit(msg)
	if err != nil {
		return
	}
	slog.Info("conventional commit", "commit", conventionalCommit)
	confRoot = filepath.Clean(confRoot)
	confPath = filepath.Clean(fmt.Sprintf("%s/%s/%s",
		confRoot, conventionalCommit.Scope, confName))
	relPath, err := filepath.Rel(confRoot, confPath)
	if err != nil {
		return
	}
	if strings.HasPrefix(relPath, "..") {
		err = fmt.Errorf("invalid scope as path: %s", conventionalCommit.Scope)
		return
	}
	groupKeywords := []string{"joj"}
	for _, groupKeyword := range groupKeywords {
		if strings.Contains(
			strings.ToLower(conventionalCommit.Description), groupKeyword) {
			group = groupKeyword
			break
		}
	}
	return
}

func ListValidScopes(confRoot, confName string) ([]string, error) {
	confRoot = filepath.Clean(confRoot)
	validScopes := []string{}
	err := filepath.Walk(confRoot, func(
		path string, info os.FileInfo, err error,
	) error {
		if err != nil {
			slog.Error("list valid scopes", "error", err)
			return err
		}
		if info.IsDir() {
			confPath := filepath.Join(path, confName)
			if _, err := os.Stat(confPath); err == nil {
				relPath, err := filepath.Rel(confRoot, path)
				if err != nil {
					return err
				}
				if relPath == "." {
					relPath = ""
				}
				validScopes = append(validScopes,
					fmt.Sprintf("'%s'", relPath))
			}
		}
		return nil
	})
	return validScopes, err
}
