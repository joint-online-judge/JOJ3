// Package keyword implements keyword-based output analysis functionality.
// It evaluates output files by searching for specific keywords and assigns scores based on matches.
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
