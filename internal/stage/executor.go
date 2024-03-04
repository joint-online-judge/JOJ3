package stage

var executorMap = map[string]Executor{}

type Executor interface {
	Run(Cmd) (*Result, error)
	Cleanup() error
}

func RegisterExecutor(name string, executor Executor) {
	executorMap[name] = executor
}
