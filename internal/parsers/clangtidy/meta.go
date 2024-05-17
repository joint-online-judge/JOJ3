package clangtidy

import "focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/internal/stage"

var name = "clangtidy"

func init() {
	stage.RegisterParser(name, &ClangTidy{})
}
