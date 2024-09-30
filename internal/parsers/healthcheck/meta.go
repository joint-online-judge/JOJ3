package healthcheck

import "focs.ji.sjtu.edu.cn/git/JOJ/JOJ3/internal/stage"

var name = "healthcheck"

func init() {
	stage.RegisterParser(name, &Healthcheck{})
}
