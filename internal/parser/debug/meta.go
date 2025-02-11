package debug

import "github.com/joint-online-judge/JOJ3/internal/stage"

var name = "debug"

type Conf struct{}

type Debug struct{}

func init() {
	stage.RegisterParser(name, &Debug{})
}
