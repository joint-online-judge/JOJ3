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
	  name = "dummy"
	  [stages.executor.with]
	    args = [ "/usr/bin/cat", "/dev/null" ]
	  [stages.parser]
	  name = "dummy"
	  [stages.parser.with]
	    score = 100
	    comment = "dummy comment for stage 0"
	[[stages]]
	name = "stage 1"
	  [stages.executor]
	  name = "dummy"
	  [stages.executor.with]
	    args = [ "/usr/bin/cat", "/dev/null" ]
	  [stages.parser]
	  name = "dummy"
	  [stages.parser.with]
	    score = 101
	    comment = "dummy comment for stage 1"
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
