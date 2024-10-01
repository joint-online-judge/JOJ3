package cppcheck

import "github.com/joint-online-judge/JOJ3/internal/stage"

var name = "cppcheck"

func init() {
	stage.RegisterParser(name, &CppCheck{})
}
