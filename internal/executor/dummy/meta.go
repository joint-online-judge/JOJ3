// Package dummy provides a mock executor implementation for testing purposes
// and serves as a template for new executor development. It always returns
// a empty accepted result.
package dummy

import "github.com/joint-online-judge/JOJ3/internal/stage"

var name = "dummy"

type Dummy struct{}

func init() {
	stage.RegisterExecutor(name, &Dummy{})
}
