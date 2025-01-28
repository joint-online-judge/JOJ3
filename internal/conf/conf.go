package conf

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

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

func ParseConventionalCommit(commit string) (*ConventionalCommit, error) {
	re := regexp.MustCompile(`(?s)^(\w+)(\(([^)]+)\))?!?: (.+?(\[([^\]]+)\])?)(\n\n(.+?))?(\n\n(.+))?$`)
	matches := re.FindStringSubmatch(strings.TrimSpace(commit))
	if matches == nil {
		return nil, fmt.Errorf("invalid conventional commit format")
	}
	cc := &ConventionalCommit{
		Type:        matches[1],
		Scope:       matches[3],
		Description: strings.TrimSpace(matches[4]),
		Group:       matches[6],
		Body:        strings.TrimSpace(matches[8]),
		Footer:      strings.TrimSpace(matches[10]),
	}
	return cc, nil
}

func ParseConfFile(path string) (conf *Conf, name string, err error) {
	conf = new(Conf)
	d := &multiconfig.DefaultLoader{}
	d.Loader = multiconfig.MultiLoader(
		&multiconfig.TagLoader{},
		&multiconfig.JSONLoader{Path: path},
	)
	d.Validator = multiconfig.MultiValidator(&multiconfig.RequiredValidator{})
	if err = d.Load(conf); err != nil {
		slog.Error("parse stages conf", "error", err)
		return
	}
	if err = d.Validate(conf); err != nil {
		slog.Error("validate stages conf", "error", err)
		return
	}
	name = conf.Name
	return
}

func GetSHA256(filePath string) (string, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Calculate SHA-256
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func parseMsg(confRoot, confName, msg, tag string) (
	confPath string, conventionalCommit *ConventionalCommit, err error,
) {
	slog.Info("parse msg", "msg", msg)
	conventionalCommit, err = ParseConventionalCommit(msg)
	if err != nil {
		return
	}
	slog.Info("conventional commit", "commit", conventionalCommit)
	confRoot = filepath.Clean(confRoot)
	confPath = filepath.Join(confRoot, conventionalCommit.Scope, confName)
	relPath, err := filepath.Rel(confRoot, confPath)
	if err != nil {
		return
	}
	if strings.HasPrefix(relPath, "..") {
		err = fmt.Errorf("invalid scope as path: %s", conventionalCommit.Scope)
		return
	}
	if tag != "" && conventionalCommit.Scope != tag {
		err = fmt.Errorf("tag does not match scope: %s != %s", tag,
			conventionalCommit.Scope)
		return
	}
	return
}

func hintValidScopes(confRoot, confName string) {
	confRoot = filepath.Clean(confRoot)
	validScopes := []string{}
	_ = filepath.Walk(confRoot, func(
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
	slog.Info("HINT: use valid scopes in commit message",
		"valid scopes", validScopes)
}

func GetConfPath(confRoot, confName, fallbackConfName, msg, tag string) (
	confPath string, confStat fs.FileInfo,
	conventionalCommit *ConventionalCommit, err error,
) {
	confPath, conventionalCommit, err = parseMsg(confRoot, confName, msg, tag)
	if err != nil {
		slog.Error("parse msg", "error", err)
		// fallback to conf file in conf root on parse error
		confPath = filepath.Join(confRoot, fallbackConfName)
		slog.Info("fallback to conf", "path", confPath)
	}
	confStat, err = os.Stat(confPath)
	if err != nil {
		if os.IsNotExist(err) {
			hintValidScopes(confRoot, confName)
		}
		slog.Error("stat conf", "error", err)
		// fallback to conf file in conf root on conf not exist
		confPath = filepath.Join(confRoot, fallbackConfName)
		slog.Info("fallback to conf", "path", confPath)
		confStat, err = os.Stat(confPath)
		if err != nil {
			slog.Error("stat fallback conf", "error", err)
			return
		}
	}
	return
}

func CheckExpire(conf *Conf) error {
	if conf.ExpireUnixTimestamp > 0 &&
		conf.ExpireUnixTimestamp < time.Now().Unix() {
		return fmt.Errorf("config file expired: %d", conf.ExpireUnixTimestamp)
	}
	return nil
}

func MatchGroups(conf *Conf, conventionalCommit *ConventionalCommit) []string {
	seen := make(map[string]bool)
	keywords := []string{}
	loweredCommitGroup := strings.ToLower(conventionalCommit.Group)
	for i, stage := range conf.Stage.Stages {
		if loweredCommitGroup == "all" {
			conf.Stage.Stages[i].Group = ""
		}
		if stage.Group == "" {
			continue
		}
		keyword := strings.ToLower(stage.Group)
		if _, exists := seen[keyword]; !exists {
			seen[keyword] = true
			keywords = append(keywords, keyword)
		}
	}
	slog.Info("group keywords from stages", "keywords", keywords)
	groups := []string{}
	for _, keyword := range keywords {
		if strings.Contains(loweredCommitGroup, keyword) {
			groups = append(groups, keyword)
		}
	}
	slog.Info("matched groups", "groups", groups)
	return groups
}
