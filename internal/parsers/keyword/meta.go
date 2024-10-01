package keyword

import "github.com/joint-online-judge/JOJ3/internal/stage"

var name = "keyword"

func init() {
	stage.RegisterParser(name, &Keyword{})
}
