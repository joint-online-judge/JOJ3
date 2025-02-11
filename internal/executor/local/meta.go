// Package local implements an executor that runs commands directly on the local
// system. It passes current environment variables to the command, which can be
// used for passing run time parameters.
package local

import "github.com/joint-online-judge/JOJ3/internal/stage"

var name = "local"

type Local struct{}

func init() {
	stage.RegisterExecutor(name, &Local{})
}
