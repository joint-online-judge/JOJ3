package main

import (
	"encoding/json"
	"flag"
	"log/slog"
	"os"

	_ "focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/internal/executors"
	_ "focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/internal/parsers"
	"focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/internal/stage"

	"github.com/pelletier/go-toml/v2"
)

func parseConfFile(tomlPath *string) (Conf, []stage.Stage) {
	tomlConfig, err := os.ReadFile(*tomlPath)
	if err != nil {
		slog.Error("read toml config", "error", err)
		os.Exit(1)
	}
	conf := Conf{
		LogLevel:   0,
		OutputPath: "joj3_result.json",
	}
	err = toml.Unmarshal(tomlConfig, &conf)
	if err != nil {
		slog.Error("parse stages config", "error", err)
		os.Exit(1)
	}
	stages := []stage.Stage{}
	for _, s := range conf.Stages {
		stages = append(stages, stage.Stage{
			Name:         s.Name,
			ExecutorName: s.Executor.Name,
			ExecutorCmd:  s.Executor.With,
			ParserName:   s.Parser.Name,
			ParserConfig: s.Parser.With,
		})
	}
	return conf, stages
}

func setupSlog(conf Conf) {
	lvl := new(slog.LevelVar)
	lvl.Set(slog.Level(conf.LogLevel))
	opts := &slog.HandlerOptions{Level: lvl}
	handler := slog.NewTextHandler(os.Stderr, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)
}

func outputResult(conf Conf, results []stage.StageResult) error {
	content, err := json.Marshal(results)
	if err != nil {
		return err
	}
	return os.WriteFile(conf.OutputPath,
		append(content, []byte("\n")...), 0o666)
}

func main() {
	tomlPath := flag.String("c", "conf.toml", "file path of the toml config")
	flag.Parse()
	conf, stages := parseConfFile(tomlPath)
	setupSlog(conf)
	defer stage.Cleanup()
	results := stage.Run(stages)
	err := outputResult(conf, results)
	if err != nil {
		slog.Error("output result", "error", err)
	}
}
