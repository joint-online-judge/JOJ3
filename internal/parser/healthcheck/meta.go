package healthcheck

import "github.com/joint-online-judge/JOJ3/internal/stage"

var name = "healthcheck"

type Healthcheck struct{}

type Conf struct {
	Stdout string `default:"stdout"`
	Stderr string `default:"stderr"`
}

func init() {
	stage.RegisterParser(name, &Healthcheck{})
}
