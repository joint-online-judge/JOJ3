// Package elf parses output of the elf static analysis tool to
// assign scores based on detected code issues.

package elf

import "github.com/joint-online-judge/JOJ3/internal/stage"

var name = "elf"

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

type Elf struct{}

func init() {
	stage.RegisterParser(name, &Elf{})
}
