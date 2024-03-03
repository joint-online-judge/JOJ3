package stage

import "github.com/criyle/go-judge/cmd/go-judge/model"

var parserMap = map[string]Parser{}

type Parser interface {
	Run(model.Result, string) ParserResult
}

type ParserResult struct {
	Score   int
	Comment string
}

func RegisterParser(name string, parser Parser) {
	parserMap[name] = parser
}
