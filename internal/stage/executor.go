package stage

var executorMap = map[string]Executor{}

type Executor interface {
	Run(Cmd) (*Result, error)
}

func RegisterExecutor(name string, executor Executor) {
	executorMap[name] = executor
}
