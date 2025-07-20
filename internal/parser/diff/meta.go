// Package diff implements string comparison functionality for the specific
// output files, comparing then with expected answers and assigning scores based
// on results.
package diff

import "github.com/joint-online-judge/JOJ3/internal/stage"

var name = "diff"

type Output struct {
	Score            int
	Filename         string
	AnswerPath       string
	CompareSpace     bool
	AlwaysHide       bool
	ForceQuitOnDiff  bool
	MaxDiffLength    int `default:"4096"` // just for reference
	MaxDiffLines     int `default:"50"`   // just for reference
	HideCommonPrefix bool
}

type Case struct {
	Outputs []Output
}

type Conf struct {
	PassComment       string `default:"🥳Passed!"`
	FailComment       string `default:"🧐Failed..."`
	FailOnNotAccepted bool   `default:"true"`
	ForceQuitOnFailed bool   `default:"false"`
	Cases             []Case
}

type Diff struct{}

func init() {
	stage.RegisterParser(name, &Diff{})
}
