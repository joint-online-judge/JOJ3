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

type Conf struct {
	SandboxExecServer string `default:"localhost:5051"`
	SandboxToken      string `default:""`
	LogPath           string `default:""`
	OutputPath        string `default:"joj3_result.json"`
	GradingRepoName   string `default:""`
	SkipTeapot        bool   `default:"true"`
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

func parseConfFile(path string) (conf Conf, err error) {
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
	return
}

func ParseMsg(confRoot, confName, msg string) (conf Conf, group string, err error) {
	slog.Info("parse msg", "msg", msg)
	conventionalCommit, err := parseConventionalCommit(msg)
	if err != nil {
		return
	}
	slog.Info("conventional commit", "commit", conventionalCommit)
	confRoot = filepath.Clean(confRoot)
	confPath := filepath.Clean(fmt.Sprintf("%s/%s/%s",
		confRoot, conventionalCommit.Scope, confName))
	relPath, err := filepath.Rel(confRoot, confPath)
	if err != nil {
		return
	}
	if strings.HasPrefix(relPath, "..") {
		err = fmt.Errorf("invalid scope as path: %s", conventionalCommit.Scope)
		return
	}
	slog.Info("try to load conf", "path", confPath)
	conf, err = parseConfFile(confPath)
	if err != nil {
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
	slog.Debug("conf loaded", "conf", conf)
	return
}

func ListValidScopes(confRoot, confName, msg string) ([]string, error) {
	conventionalCommit, err := parseConventionalCommit(msg)
	if err != nil {
		return []string{}, err
	}
	slog.Info("conventional commit", "commit", conventionalCommit)
	confRoot = filepath.Clean(confRoot)
	validScopes := []string{}
	err = filepath.Walk(confRoot, func(path string, info os.FileInfo, err error) error {
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
				validScopes = append(validScopes, relPath)
			}
		}
		return nil
	})
	return validScopes, err
}
