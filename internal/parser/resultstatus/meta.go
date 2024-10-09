package resultstatus

import "github.com/joint-online-judge/JOJ3/internal/stage"

var name = "result-status"

func init() {
	stage.RegisterParser(name, &ResultStatus{})
}
