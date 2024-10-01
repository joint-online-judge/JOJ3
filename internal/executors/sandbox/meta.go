package sandbox

import (
	"github.com/joint-online-judge/JOJ3/internal/stage"
)

var name = "sandbox"

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
