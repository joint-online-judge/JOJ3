package plugin

import (
	"fmt"
	"plugin"

	"github.com/joint-online-judge/JOJ3/internal/stage"
)

func (*Plugin) Run(results []stage.ExecutorResult, confAny any) (
	[]stage.ParserResult, bool, error,
) {
	conf, err := stage.DecodeConf[Conf](confAny)
	if err != nil {
		return nil, true, err
	}
	plug, err := plugin.Open(conf.ModPath)
	if err != nil {
		return nil, true, err
	}
	symParser, err := plug.Lookup(conf.SymName)
	if err != nil {
		return nil, true, err
	}
	var parser stage.Parser
	parser, ok := symParser.(stage.Parser)
	if !ok {
		return nil, true, fmt.Errorf("unexpected type from module symbol")
	}
	return parser.Run(results, confAny)
}
