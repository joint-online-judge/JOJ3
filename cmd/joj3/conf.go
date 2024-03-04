package main

import "focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/internal/stage"

type Conf struct {
	LogLevel   int
	OutputPath string
	Stages     []struct {
		Name     string
		Executor struct {
			Name string
			With stage.Cmd
		}
		Parser struct {
			Name string
			With interface{}
		}
	}
}
