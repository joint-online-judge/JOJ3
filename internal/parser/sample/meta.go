package sample

import "github.com/joint-online-judge/JOJ3/internal/stage"

var name = "sample"

type Conf struct {
	Score   int
	Comment string
	Stdout  string `default:"stdout"`
	Stderr  string `default:"stderr"`
}

type Sample struct{}

func init() {
	stage.RegisterParser(name, &Sample{})
}
