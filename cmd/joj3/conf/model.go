package conf

import (
	"github.com/joint-online-judge/JOJ3/internal/stage"
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
	Name                   string `default:"unknown"`
	LogPath                string `default:""`
	ActorCsvPath           string `default:""`
	EffectiveUnixTimestamp int64  `default:"-1"`
	ExpireUnixTimestamp    int64  `default:"-1"`
	MaxTotalScore          int    `default:"-1"`
	Stage                  struct {
		SandboxExecServer string `default:"localhost:5051"`
		SandboxToken      string `default:""`
		OutputPath        string `default:"joj3_result.json"`
		PreStages         []ConfStage
		Stages            []ConfStage
		PostStages        []ConfStage
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
	Group       string
	Body        string
	Footer      string
}
