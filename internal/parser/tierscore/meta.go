// Package tierscore provides a parser for tiered scoring based on
// time and memory constraints. Leave the field empty or 0 to disable.
package tierscore

import "github.com/joint-online-judge/JOJ3/internal/stage"

var name = "tierscore"

type Tier struct {
	TimeLessThan   uint64 // ns
	MemoryLessThan uint64 // bytes
	Score          int
}

type Conf struct {
	Tiers []Tier
}

type TierScore struct{}

func init() {
	stage.RegisterParser(name, &TierScore{})
}
