package diff

import "focs.ji.sjtu.edu.cn/git/JOJ/JOJ3/internal/stage"

var name = "diff"

func init() {
	stage.RegisterParser(name, &Diff{})
}
