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

type ConfGroup struct {
	Name           string
	MaxCount       int
	TimePeriodHour int
}

type Conf struct {
	Name                string `default:"unknown"`
	LogPath             string `default:""`
	ActorCsvPath        string `default:""`
	ExpireUnixTimestamp int64  `default:"-1"`
	MaxTotalScore       int    `default:"-1"`
	Stage               struct {
		SandboxExecServer string `default:"localhost:5051"`
		SandboxToken      string `default:""`
		OutputPath        string `default:"joj3_result.json"`
		Stages            []ConfStage
	}
	Teapot struct {
		LogPath               string `default:"/home/tt/.cache/joint-teapot-debug.log"`
		EnvFilePath           string `default:"/home/tt/.config/teapot/teapot.env"`
		ScoreboardPath        string `default:"scoreboard.csv"`
		FailedTablePath       string `default:"failed-table.md"`
		GradingRepoName       string `default:""`
		SkipIssue             bool   `default:"false"`
		SkipScoreboard        bool   `default:"false"`
		SkipFailedTable       bool   `default:"false"`
		SubmitterInIssueTitle bool   `default:"true"`
		Groups                []ConfGroup
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
