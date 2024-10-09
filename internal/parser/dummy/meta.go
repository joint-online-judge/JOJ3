package dummy

import "github.com/joint-online-judge/JOJ3/internal/stage"

var name = "dummy"

func init() {
	stage.RegisterParser(name, &Dummy{})
}
