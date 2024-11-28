package executors

import (
	_ "github.com/joint-online-judge/JOJ3/internal/executor/dummy"
	_ "github.com/joint-online-judge/JOJ3/internal/executor/local"
	"github.com/joint-online-judge/JOJ3/internal/executor/sandbox"
)

// this file does nothing but imports to ensure all the init() functions
// in the subpackages are called

// overwrite the default registered executors
func InitWithConf(sandboxExecServer, sandboxToken string) {
	sandbox.InitWithConf(sandboxExecServer, sandboxToken)
}
