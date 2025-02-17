package parsers

import (
	_ "github.com/joint-online-judge/JOJ3/internal/parser/clangtidy"
	_ "github.com/joint-online-judge/JOJ3/internal/parser/cppcheck"
	_ "github.com/joint-online-judge/JOJ3/internal/parser/cpplint"
	_ "github.com/joint-online-judge/JOJ3/internal/parser/debug"
	_ "github.com/joint-online-judge/JOJ3/internal/parser/diff"
	_ "github.com/joint-online-judge/JOJ3/internal/parser/dummy"
	_ "github.com/joint-online-judge/JOJ3/internal/parser/healthcheck"
	_ "github.com/joint-online-judge/JOJ3/internal/parser/keyword"
	_ "github.com/joint-online-judge/JOJ3/internal/parser/log"
	_ "github.com/joint-online-judge/JOJ3/internal/parser/plugin"
	_ "github.com/joint-online-judge/JOJ3/internal/parser/resultdetail"
	_ "github.com/joint-online-judge/JOJ3/internal/parser/resultstatus"
	_ "github.com/joint-online-judge/JOJ3/internal/parser/sample"
	_ "github.com/joint-online-judge/JOJ3/internal/parser/tierscore"
)

// this file does nothing but imports to ensure all the init() functions
// in the subpackages are called
