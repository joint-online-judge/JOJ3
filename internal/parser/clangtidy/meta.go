package clangtidy

import "github.com/joint-online-judge/JOJ3/internal/stage"

var name = "clangtidy"

func init() {
	stage.RegisterParser(name, &ClangTidy{})
}
