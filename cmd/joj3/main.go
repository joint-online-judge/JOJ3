package main

import (
	"fmt"

	_ "focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/internal/executors"
	_ "focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/internal/parsers"
	"focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/internal/stage"
)

func main() {
	stages := stage.ParseStages()
	results := stage.Run(stages)
	for _, result := range results {
		fmt.Printf(
			"%s: score: %d, comment: %s\n",
			result.Name, result.Score, result.Comment,
		)
	}
}
