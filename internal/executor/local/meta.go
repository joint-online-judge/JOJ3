package local

import "github.com/joint-online-judge/JOJ3/internal/stage"

var name = "local"

func init() {
	stage.RegisterExecutor(name, &Local{})
}
