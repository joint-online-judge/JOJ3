package resultstatus

import "focs.ji.sjtu.edu.cn/git/JOJ/JOJ3/internal/stage"

var name = "result-status"

func init() {
	stage.RegisterParser(name, &ResultStatus{})
}
