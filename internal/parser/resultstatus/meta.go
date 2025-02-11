// Package resultstatus provides functionality to parse execution results
// and determine success/failure status. It can return forced quit status
// when a non-accepted status is encountered.
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
