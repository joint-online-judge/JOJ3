package file

import (
	"strings"

	"github.com/joint-online-judge/JOJ3/internal/stage"
)

type Conf struct {
	Name                string
	ForceQuitOnNonEmpty bool
}

type File struct{}

func (*File) Run(results []stage.ExecutorResult, confAny any) (
	[]stage.ParserResult, bool, error,
) {
	conf, err := stage.DecodeConf[Conf](confAny)
	if err != nil {
		return nil, true, err
	}
	var res []stage.ParserResult
	forceQuit := false
	for _, result := range results {
		content := result.Files[conf.Name]
		if conf.ForceQuitOnNonEmpty && content != "" {
			forceQuit = true
		}
		if !strings.HasSuffix(content, "\n") {
			content += "\n"
		}
		res = append(res, stage.ParserResult{Score: 0, Comment: content})
	}
	return res, forceQuit, nil
}
