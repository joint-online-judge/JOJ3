package executors

import (
	_ "focs.ji.sjtu.edu.cn/git/JOJ/JOJ3/internal/executors/dummy"
	"focs.ji.sjtu.edu.cn/git/JOJ/JOJ3/internal/executors/sandbox"
)

// this file does nothing but imports to ensure all the init() functions
// in the subpackages are called

// overwrite the default registered executors
func InitWithConf(sandboxExecServer, sandboxToken string) {
	sandbox.InitWithConf(sandboxExecServer, sandboxToken)
}
