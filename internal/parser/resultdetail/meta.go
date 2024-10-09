package resultdetail

import "github.com/joint-online-judge/JOJ3/internal/stage"

var name = "result-detail"

func init() {
	stage.RegisterParser(name, &ResultDetail{})
}
