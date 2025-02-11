// Package clangtidy parses output of the cppcheck static analysis tool to
// assign scores based on detected code issues.
// Check examples on running cppcheck for parseable output.
package cppcheck

import "github.com/joint-online-judge/JOJ3/internal/stage"

var name = "cppcheck"

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

type CppCheck struct{}

func init() {
	stage.RegisterParser(name, &CppCheck{})
}
