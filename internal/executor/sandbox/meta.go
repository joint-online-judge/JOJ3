// Package sandbox provides a sandboxed execution environment for running
// untrusted code. It integrates with the go-judge execution service to provide
// isolated and secure code execution. By default, it uses gRPC to communicate
// with go-judge.
package sandbox

import (
	"github.com/criyle/go-judge/pb"
	"github.com/joint-online-judge/JOJ3/internal/stage"
)

var name = "sandbox"

type Sandbox struct {
	execServer, token string
	cachedMap         map[string]string
	execClient        pb.ExecutorClient
}

func init() {
	stage.RegisterExecutor(name, &Sandbox{
		execServer: "localhost:5051",
		token:      "",
		cachedMap:  make(map[string]string),
	})
}

// overwrite the default registered executor
func InitWithConf(execServer, token string) {
	stage.RegisterExecutor(name, &Sandbox{
		execServer: execServer,
		token:      token,
		cachedMap:  make(map[string]string),
	})
}
