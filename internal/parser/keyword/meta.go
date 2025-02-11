package keyword

import "github.com/joint-online-judge/JOJ3/internal/stage"

var name = "keyword"

type Match struct {
	Keywords      []string
	Score         int
	MaxMatchCount int
}

type Conf struct {
	Score             int
	Files             []string
	ForceQuitOnDeduct bool `default:"false"`
	Matches           []Match
}

type Keyword struct{}

func init() {
	stage.RegisterParser(name, &Keyword{})
}
