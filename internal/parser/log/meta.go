// Package log logs the json key-value pairs from given file. The log can be
// used for Loki that contains run time status.
package log

import "github.com/joint-online-judge/JOJ3/internal/stage"

var name = "log"

type Conf struct {
	FileName string `default:"stdout"`
	Msg      string `default:"log msg"`
	Level    int    `default:"0"`
}

type Log struct{}

func init() {
	stage.RegisterParser(name, &Log{})
}
