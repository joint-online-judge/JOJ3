package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"regexp"

	"focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/internal/stage"
	"github.com/go-git/go-git/v5"
	"github.com/koding/multiconfig"
)

type Conf struct {
	SandboxExecServer string `default:"localhost:5051"`
	SandboxToken      string `default:""`
	LogLevel          int    `default:"0"`
	OutputPath        string `default:"joj3_result.json"`
	Stages            []struct {
		Name     string
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

// TODO: add other fields to match? not only limit to latest commit message
type MetaConf struct {
	Patterns []struct {
		Filename string
		Regex    string
	}
}

func parseMetaConfFile(path string) (metaConf MetaConf, err error) {
	// FIXME: remove this default meta config, it is only for demonstration
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return MetaConf{
			Patterns: []struct {
				Filename string
				Regex    string
			}{
				{
					Filename: "conf.json",
					Regex:    ".*",
				},
			},
		}, nil
	}
	d := &multiconfig.DefaultLoader{}
	d.Loader = multiconfig.MultiLoader(
		&multiconfig.TagLoader{},
		&multiconfig.JSONLoader{Path: path},
	)
	d.Validator = multiconfig.MultiValidator(&multiconfig.RequiredValidator{})
	if err = d.Load(&metaConf); err != nil {
		slog.Error("parse meta conf", "error", err)
		return
	}
	if err = d.Validate(&metaConf); err != nil {
		slog.Error("validate meta conf", "error", err)
		return
	}
	return
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

func commitMsgToConf(metaConfPath string) (conf Conf, err error) {
	metaConf, err := parseMetaConfFile(metaConfPath)
	if err != nil {
		return
	}
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
	msg := commit.Message
	slog.Debug("commit msg to conf", "msg", msg)
	for _, pattern := range metaConf.Patterns {
		if matched, _ := regexp.MatchString(pattern.Regex, msg); matched {
			slog.Debug("pattern matched", "pattern", pattern)
			return parseConfFile(pattern.Filename)
		}
	}
	err = fmt.Errorf("no pattern matched")
	return
}
