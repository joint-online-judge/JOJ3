// Package plugin provides functionality to load and run parser plugins
// dynamically. It is used for custom parsers.
// The plugin needs to be located at `ModPath` and export a symbol with name
// `SymName` that implements the stage.Parser interface.
package plugin

import "github.com/joint-online-judge/JOJ3/internal/stage"

var name = "plugin"

type Conf struct {
	ModPath string
	SymName string `default:"Parser"`
}

type Plugin struct{}

func init() {
	stage.RegisterParser(name, &Plugin{})
}
