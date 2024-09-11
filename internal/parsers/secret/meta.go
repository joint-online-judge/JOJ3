package secret

import "focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/internal/stage"

var name = "secret"

func init() {
	stage.RegisterParser(name, &Secret{})
}
