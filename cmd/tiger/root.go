package main

import (
	"fmt"
	"os"

	"focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/pkg/runner"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: " + os.Args[0] + " <command>")
		os.Exit(1)
	}
	// hard limit of timeout 1000ms
	runResult, err := runner.RunInCgroupsV1(os.Args[1:], "nobody", "/joj3.tiger", 1000)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	fmt.Printf("ReturnCode: %d\n", runResult.ReturnCode)
	fmt.Printf("Stdout: %s\n", runResult.Stdout)
	fmt.Printf("Stderr: %s\n", runResult.Stderr)
	fmt.Printf("TimedOut: %v\n", runResult.TimedOut)
	fmt.Printf("TimeNs: %v\n", runResult.TimeNs)
	fmt.Printf("MemoryByte: %v\n", runResult.MemoryByte)
}
