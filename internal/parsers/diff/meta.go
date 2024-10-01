package diff

import "github.com/joint-online-judge/JOJ3/internal/stage"

var name = "diff"

func init() {
	stage.RegisterParser(name, &Diff{})
}
