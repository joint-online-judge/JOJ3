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
name = "compile"
	[stages.executor]
	name = "sandbox"
	[stages.executor.with]
	args = [ "/usr/bin/g++", "a.cc", "-o", "a" ]
	env = [ "PATH=/usr/bin:/bin" ]
	cpuLimit = 10_000_000_000
	memoryLimit = 104_857_600
	procLimit = 50
	copyOut = [ "stdout", "stderr" ]
	copyOutCached = [ "a" ]
[stages.executor.with.copyIn."a.cc"]
content = """
#include <iostream>
int main() {
	int a, b;
	std::cin >> a >> b;
	std::cout << a + b << '\\n';
}"""
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
	comment = "compile done"
[[stages]]
name = "run"
	[stages.executor]
	name = "sandbox"
	[stages.executor.with]
	args = [ "a" ]
	env = [ "PATH=/usr/bin:/bin" ]
	cpuLimit = 10_000_000_000
	memoryLimit = 104_857_600
	procLimit = 50
	copyOut = [ "stdout", "stderr" ]
	[stages.executor.with.copyInCached]
	    a = "a"
		[[stages.executor.with.files]]
		content = "1 1"
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
	comment = "run done"
	`
	defer stage.Cleanup()
	stages := stage.ParseStages(tomlConfig)
	results := stage.Run(stages)
	for _, result := range results {
		fmt.Printf(
			"%s: score: %d, comment: %s\n",
			result.Name, result.Score, result.Comment,
		)
	}
}
