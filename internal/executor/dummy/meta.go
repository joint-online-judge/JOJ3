package dummy

import "github.com/joint-online-judge/JOJ3/internal/stage"

var name = "dummy"

type Dummy struct{}

func init() {
	stage.RegisterExecutor(name, &Dummy{})
}
