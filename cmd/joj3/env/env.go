package env

import (
	"fmt"
	"os"
	"time"
)

type Attribute struct {
	ConfName   string
	Groups     string
	RunID      string
	Actor      string
	Repository string
	Sha        string
	Ref        string
	Workflow   string
	RunNumber  string
	ActorName  string
	ActorID    string
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
	os.Setenv("CONF_NAME", Attr.ConfName)
	os.Setenv("GROUPS", Attr.Groups)
	os.Setenv("RUN_ID", Attr.RunID)
}
