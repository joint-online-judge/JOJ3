package dummy

import (
	"fmt"
	"log/slog"
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
		slog.Error("dummy negative score", "score", conf.Score)
		err = fmt.Errorf("dummy negative score: %d", conf.Score)
		return
	}
	res.Score = conf.Score
	res.Comment = "dummy comment"
	return
}
