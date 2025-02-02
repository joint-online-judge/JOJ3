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
)

type Attribute struct {
	ConfName           string
	CommitMsg          string
	Groups             string
	RunID              string
	Actor              string
	Repository         string
	Sha                string
	Ref                string
	Workflow           string
	RunNumber          string
	ActorName          string
	ActorID            string
	ForceQuitStageName string
	OutputPath         string
}

var Attr Attribute

func init() {
	timestamp := time.Now().UnixNano()
	pid := os.Getpid()
	high := timestamp >> 32
	low := timestamp & 0xFFFFFFFF
	combined := high ^ low
	combined ^= int64(pid)
	combined ^= int64(timestamp >> 16)
	combined ^= (combined >> 8)
	combined ^= (combined << 16)
	Attr.RunID = fmt.Sprintf("%08X", combined&0xFFFFFFFF)
	Attr.Actor = os.Getenv("GITHUB_ACTOR")
	Attr.Repository = os.Getenv("GITHUB_REPOSITORY")
	Attr.Sha = os.Getenv("GITHUB_SHA")
	Attr.Ref = os.Getenv("GITHUB_REF")
	Attr.Workflow = os.Getenv("GITHUB_WORKFLOW")
	Attr.RunNumber = os.Getenv("GITHUB_RUN_NUMBER")
}

func Set() {
	os.Setenv(ConfName, Attr.ConfName)
	os.Setenv(Groups, Attr.Groups)
	os.Setenv(RunID, Attr.RunID)
	os.Setenv(CommitMsg, Attr.CommitMsg)
	os.Setenv(ForceQuitStageName, Attr.ForceQuitStageName)
	os.Setenv(OutputPath, Attr.OutputPath)
}
