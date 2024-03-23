package sandbox

import (
	"focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/internal/stage"
)

var name = "sandbox"

func init() {
	stage.RegisterExecutor(name, &Sandbox{
		// TODO: read from conf
		execServer: "localhost:5051",
		token:      "",
		cachedMap:  make(map[string]string),
	})
}
