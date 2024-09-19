package main

import (
	"log/slog"

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
	CopyInCwd    *bool

	CopyOut       *[]string
	CopyOutCached *[]string
	CopyOutMax    *uint64
	CopyOutDir    *string

	TTY               *bool
	StrictMemoryLimit *bool
	DataSegmentLimit  *bool
	AddressSpaceLimit *bool
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
	slog.Error("parse stages conf", "conf", conf)
	return
}

func commitMsgToConf() (conf Conf, err error) {
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
	// TODO: parse msg to conf name
	conf, err = parseConfFile("conf.json")
	return
}
