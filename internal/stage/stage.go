package stage

import "encoding/json"

type StageExecutor struct {
	Name string
	Cmds []Cmd
}
type StageParser struct {
	Name string
	Conf any
}

type Stage struct {
	Name     string
	Executor StageExecutor
	Parsers  []StageParser
}

type NonNullSlice[T any] []T

func (s NonNullSlice[T]) MarshalJSON() ([]byte, error) {
	if len(s) == 0 {
		return []byte("[]"), nil
	}
	return json.Marshal([]T(s))
}

type StageResult struct {
	Name      string                     `json:"name"`
	Results   NonNullSlice[ParserResult] `json:"results"`
	ForceQuit bool                       `json:"force_quit"` // underscore as it will dump to file
}

type CaseDetail struct {
	Index          int            `json:"index"`
	ExecutorResult ExecutorResult `json:"executorResult"`
	ParserScores   map[string]int `json:"parserScores"`
}

type StageDetail struct {
	Name        string       `json:"name"`
	CaseDetails []CaseDetail `json:"caseDetails"`
}
