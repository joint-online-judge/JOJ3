// Package sample provides functionality to parse and process sample outputs
// from stdout and stderr of the sample program. Use this as a sample.
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
