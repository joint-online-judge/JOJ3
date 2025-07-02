package env

import (
	"fmt"
	"os"
	"time"
)

const (
	ConfName           = "JOJ3_CONF_NAME"
	Groups             = "JOJ3_GROUPS"
	RunID              = "JOJ3_RUN_ID"
	CommitMsg          = "JOJ3_COMMIT_MSG"
	ForceQuitStageName = "JOJ3_FORCE_QUIT_STAGE_NAME"
	OutputPath         = "JOJ3_OUTPUT_PATH"

	GitHubActor      = "GITHUB_ACTOR"
	GitHubRepository = "GITHUB_REPOSITORY"
	GitHubSha        = "GITHUB_SHA"
	GitHubRef        = "GITHUB_REF"
	GitHubWorkflow   = "GITHUB_WORKFLOW"
	GitHubRunNumber  = "GITHUB_RUN_NUMBER"
)

var runID string

func generateRunID() string {
	timestamp := time.Now().UnixNano()
	pid := os.Getpid()
	high := timestamp >> 32
	low := timestamp & 0xFFFFFFFF
	combined := high ^ low
	combined ^= int64(pid)
	combined ^= timestamp >> 16
	combined ^= (combined >> 8)
	combined ^= (combined << 16)
	return fmt.Sprintf("%08X", combined&0xFFFFFFFF)
}

func init() {
	runID = generateRunID()
}

func GetRunID() string              { return runID }
func GetConfName() string           { return os.Getenv(ConfName) }
func GetGroups() string             { return os.Getenv(Groups) }
func GetCommitMsg() string          { return os.Getenv(CommitMsg) }
func GetForceQuitStageName() string { return os.Getenv(ForceQuitStageName) }
func GetOutputPath() string         { return os.Getenv(OutputPath) }

func SetConfName(val string)           { os.Setenv(ConfName, val) }
func SetGroups(val string)             { os.Setenv(Groups, val) }
func SetCommitMsg(val string)          { os.Setenv(CommitMsg, val) }
func SetForceQuitStageName(val string) { os.Setenv(ForceQuitStageName, val) }
func SetOutputPath(val string)         { os.Setenv(OutputPath, val) }

func GetActor() string      { return os.Getenv(GitHubActor) }
func GetRepository() string { return os.Getenv(GitHubRepository) }
func GetSha() string        { return os.Getenv(GitHubSha) }
func GetRef() string        { return os.Getenv(GitHubRef) }
func GetWorkflow() string   { return os.Getenv(GitHubWorkflow) }
func GetRunNumber() string  { return os.Getenv(GitHubRunNumber) }
