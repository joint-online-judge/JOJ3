package parsers

import (
	_ "github.com/joint-online-judge/JOJ3/internal/parsers/clangtidy"
	_ "github.com/joint-online-judge/JOJ3/internal/parsers/cppcheck"
	_ "github.com/joint-online-judge/JOJ3/internal/parsers/cpplint"
	_ "github.com/joint-online-judge/JOJ3/internal/parsers/diff"
	_ "github.com/joint-online-judge/JOJ3/internal/parsers/dummy"
	_ "github.com/joint-online-judge/JOJ3/internal/parsers/healthcheck"
	_ "github.com/joint-online-judge/JOJ3/internal/parsers/keyword"
	_ "github.com/joint-online-judge/JOJ3/internal/parsers/resultstatus"
	_ "github.com/joint-online-judge/JOJ3/internal/parsers/sample"
)

// this file does nothing but imports to ensure all the init() functions
// in the subpackages are called
