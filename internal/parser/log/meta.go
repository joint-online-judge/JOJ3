package log

import "github.com/joint-online-judge/JOJ3/internal/stage"

var name = "log"

func init() {
	stage.RegisterParser(name, &Log{})
}
