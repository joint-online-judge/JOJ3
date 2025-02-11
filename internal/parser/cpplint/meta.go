// Package clangtidy parses output of the cpplint style checker tool to assign
// scores based on detected code issues.
package cpplint

import "github.com/joint-online-judge/JOJ3/internal/stage"

var name = "cpplint"

type Match struct {
	Keywords []string
	Score    int
}

type Conf struct {
	Score             int
	Matches           []Match
	Stdout            string `default:"stdout"`
	Stderr            string `default:"stderr"`
	ForceQuitOnDeduct bool   `default:"false"`
}

type Cpplint struct{}

func init() {
	stage.RegisterParser(name, &Cpplint{})
}
