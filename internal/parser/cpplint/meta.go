package cpplint

import "github.com/joint-online-judge/JOJ3/internal/stage"

var name = "cpplint"

func init() {
	stage.RegisterParser(name, &Cpplint{})
}
