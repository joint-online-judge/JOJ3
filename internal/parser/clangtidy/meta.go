package clangtidy

import "github.com/joint-online-judge/JOJ3/internal/stage"

var name = "clangtidy"

type Match struct {
	Keywords []string
	Score    int
}

type Conf struct {
	Score             int
	RootDir           string `default:"/w"`
	Matches           []Match
	Stdout            string `default:"stdout"`
	Stderr            string `default:"stderr"`
	ForceQuitOnDeduct bool   `default:"false"`
}

type ClangTidy struct{}

func init() {
	stage.RegisterParser(name, &ClangTidy{})
}
