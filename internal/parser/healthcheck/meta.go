// Package healthcheck parses the output of the repo-health-checker tool and
// return forced quit status on error.
package healthcheck

import "github.com/joint-online-judge/JOJ3/internal/stage"

var name = "healthcheck"

type Healthcheck struct{}

type Conf struct {
	Score  int    `default:"0"`
	Stdout string `default:"stdout"`
	Stderr string `default:"stderr"`
}

func init() {
	stage.RegisterParser(name, &Healthcheck{})
}
