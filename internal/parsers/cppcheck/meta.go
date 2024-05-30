package cppcheck

import "focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/internal/stage"

var name = "cppcheck"

func init() {
	stage.RegisterParser(name, &CppCheck{})
}
