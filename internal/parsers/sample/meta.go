package sample

import "focs.ji.sjtu.edu.cn/git/JOJ/JOJ3/internal/stage"

var name = "sample"

func init() {
	stage.RegisterParser(name, &Sample{})
}
