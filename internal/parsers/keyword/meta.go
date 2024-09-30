package keyword

import "focs.ji.sjtu.edu.cn/git/JOJ/JOJ3/internal/stage"

var name = "keyword"

func init() {
	stage.RegisterParser(name, &Keyword{})
}
