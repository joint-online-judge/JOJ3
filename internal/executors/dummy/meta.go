package dummy

import "focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/internal/stage"

var name = "dummy"

func init() {
	stage.RegisterExecutor(name, &Dummy{})
}
