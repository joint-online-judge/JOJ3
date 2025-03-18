package stage

var parserMap = map[string]Parser{}

type Parser interface {
	Run([]ExecutorResult, any) ([]ParserResult, bool, error)
}

func RegisterParser(name string, parser Parser) {
	parserMap[name] = parser
}

type ParserResult struct {
	Score   int    `json:"score"`
	Comment string `json:"comment"`
}
