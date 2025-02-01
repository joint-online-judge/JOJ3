package debug

import "github.com/joint-online-judge/JOJ3/internal/stage"

var name = "debug"

func init() {
	stage.RegisterParser(name, &Debug{})
}
