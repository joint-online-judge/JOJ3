package sample

import (
	"fmt"
)

type Conf struct {
	Score int
}

type Result struct {
	Score   int
	Comment string
}

func Run(conf Conf) (res Result, err error) {
	if conf.Score < 0 {
		// Just return the error here instead of logging, as it is run inside
		// the sandbox, the logs will not show in drone output directly.
		// If there are more kinds of errors need to be handled separately, add
		// more fields in the Result struct, don't mess everything up in Stderr.
		err = fmt.Errorf("sample negative score: %d", conf.Score)
		return
	}
	res.Score = conf.Score
	res.Comment = "sample comment"
	return
}
