package main

import (
	"fmt"

	_ "focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/internal/executors"
	_ "focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/internal/parsers"
	"focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/internal/stage"
)

func main() {
	tomlConfig := `
[[stages]]
name = "stage 0"
	[stages.executor]
	name = "sandbox"
	[stages.executor.with]
	args = [ "ls" ]
	env = [ "PATH=/usr/bin:/bin" ]
	cpuLimit = 10_000_000_000
	memoryLimit = 104_857_600
	procLimit = 50
	copyOut = [ "stdout", "stderr" ]
		[[stages.executor.with.files]]
		content = ""
		[[stages.executor.with.files]]
		name = "stdout"
		max = 4_096
		[[stages.executor.with.files]]
		name = "stderr"
		max = 4_096
	[stages.parser]
	name = "dummy"
	[stages.parser.with]
	score = 100
	comment = "dummy comment for stage 0"
	`
	stages := stage.ParseStages(tomlConfig)
	results := stage.Run(stages)
	for _, result := range results {
		fmt.Printf(
			"%s: score: %d, comment: %s\n",
			result.Name, result.Score, result.Comment,
		)
	}
}
