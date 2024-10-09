package sample

import "github.com/joint-online-judge/JOJ3/internal/stage"

var name = "sample"

func init() {
	stage.RegisterParser(name, &Sample{})
}
