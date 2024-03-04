package stage

var parserMap = map[string]Parser{}

type Parser interface {
	Run(*Result, any) (*ParserResult, error)
}

func RegisterParser(name string, parser Parser) {
	parserMap[name] = parser
}
