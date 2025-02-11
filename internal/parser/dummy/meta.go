package dummy

import "github.com/joint-online-judge/JOJ3/internal/stage"

var name = "dummy"

type Conf struct {
	Score     int
	Comment   string
	ForceQuit bool
}

type Dummy struct{}

func init() {
	stage.RegisterParser(name, &Dummy{})
}
