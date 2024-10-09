package healthcheck

import "github.com/joint-online-judge/JOJ3/internal/stage"

var name = "healthcheck"

func init() {
	stage.RegisterParser(name, &Healthcheck{})
}
