package resultstatus

import "github.com/joint-online-judge/JOJ3/internal/stage"

var name = "result-status"

type Conf struct {
	Score                  int
	Comment                string
	ForceQuitOnNotAccepted bool `default:"true"`
}

type ResultStatus struct{}

func init() {
	stage.RegisterParser(name, &ResultStatus{})
}
