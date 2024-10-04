package elf

import "github.com/joint-online-judge/JOJ3/internal/stage"

var name = "elf"

func init() {
	stage.RegisterParser(name, &Elf{})
}
