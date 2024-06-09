package elf

import "focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/internal/stage"

var name = "elf"

func init() {
	stage.RegisterParser(name, &Elf{})
}
