package main

import (
	"flag"
	"log/slog"
	"os"

	_ "focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/internal/executors"
	_ "focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/internal/parsers"
	"focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/internal/stage"
)

func main() {
	tomlPath := flag.String("c", "conf.toml", "file path of the toml config")
	flag.Parse()
	tomlConfig, err := os.ReadFile(*tomlPath)
	if err != nil {
		slog.Error("read toml config", "error", err)
		os.Exit(1)
	}
	defer stage.Cleanup()
	stages := stage.ParseStages(tomlConfig)
	results := stage.Run(stages)
	for _, result := range results {
		slog.Info(
			"stage result",
			"name", result.Name,
			"score", result.Score,
			"comment", result.Comment,
		)
	}
}
