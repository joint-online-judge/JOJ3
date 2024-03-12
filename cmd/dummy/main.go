package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/pkg/dummy"
)

func main() {
	score := flag.Int("score", 0, "score")
	flag.Parse()
	res, err := dummy.Run(dummy.Conf{Score: *score})
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	b, err := json.Marshal(res)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("%s", b)
}
