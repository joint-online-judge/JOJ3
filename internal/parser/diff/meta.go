// Package diff implements string comparison functionality for the specific
// output files, comparing then with expected answers and assigning scores based
// on results.
package diff

import "github.com/joint-online-judge/JOJ3/internal/stage"

var name = "diff"

type Conf struct {
	PassComment       string `default:"🥳Passed!"`
	FailComment       string `default:"🧐Failed..."`
	FailOnNotAccepted bool   `default:"true"`
	ForceQuitOnFailed bool   `default:"false"`
	Cases             []struct {
		Outputs []struct {
			Score            int
			FileName         string
			AnswerPath       string
			CompareSpace     bool
			AlwaysHide       bool
			ForceQuitOnDiff  bool
			MaxDiffLength    int `default:"2048"` // just for reference
			MaxDiffLines     int `default:"50"`   // just for reference
			HideCommonPrefix bool
		}
	}
}

type Diff struct{}

func init() {
	stage.RegisterParser(name, &Diff{})
}
