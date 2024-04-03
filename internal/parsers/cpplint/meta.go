package cpplint

import "focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/internal/stage"

var name = "cpplint"

func init() {
	stage.RegisterParser(name, &Cpplint{})
}
