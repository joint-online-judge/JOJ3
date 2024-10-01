package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/joint-online-judge/JOJ3/pkg/sample"
)

func main() {
	score := flag.Int("score", 0, "score")
	flag.Parse()
	res, err := sample.Run(sample.Conf{Score: *score})
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
