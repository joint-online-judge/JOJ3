package file

import "github.com/joint-online-judge/JOJ3/internal/stage"

var name = "file"

func init() {
	stage.RegisterParser(name, &File{})
}
