// Package main provides a repo-health-checker executable that checks the
// health of a repository. Its output should be parsed by the healthcheck
// parser.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/joint-online-judge/JOJ3/pkg/healthcheck"
)

// parseMultiValueFlag parses a multi-value command-line flag and appends its values to the provided slice.
// It registers a flag with the specified name and description, associating it with a multiStringValue receiver.
func parseMultiValueFlag(flags *flag.FlagSet, values *[]string, flagName, description string) {
	flags.Var((*multiStringValue)(values), flagName, description)
}

type multiStringValue []string

// Set appends a new value to the multiStringValue slice.
// It satisfies the flag.Value interface, allowing multiStringValue to be used as a flag value.
func (m *multiStringValue) Set(value string) error {
	*m = append(*m, value)
	return nil
}

func (m *multiStringValue) String() string {
	return fmt.Sprintf("%v", *m)
}

func setupSlog() {
	opts := &slog.HandlerOptions{}
	handler := slog.NewTextHandler(os.Stderr, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)
}

var Version string

func run(args []string, stdout io.Writer) error {
	flags := flag.NewFlagSet("repo-health-checker", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	var (
		rootDir           string
		repoSize          float64
		checkFileNameList string
		checkFileSumList  string
		metaFile          []string
		whitelistedChars  string
		allowedDomainList string
		actorCsvPath      string
		showVersion       bool
	)
	flags.BoolVar(&showVersion, "version", false, "print current version")
	flags.StringVar(&rootDir, "root", ".", "root dir for forbidden files check")
	flags.Float64Var(&repoSize, "repoSize", 2, "maximum size of the repo in MiB")
	flags.StringVar(&checkFileNameList, "checkFileNameList", "", "comma-separated list of files to check")
	flags.StringVar(&checkFileSumList, "checkFileSumList", "", "comma-separated list of expected checksums")
	flags.StringVar(&whitelistedChars, "whitelistedChars", "", "comma-separated list of non-ASCII characters allowed in files")
	flags.StringVar(&allowedDomainList, "allowedDomainList", "sjtu.edu.cn", "comma-separated list of allowed domains for commit author email")
	flags.StringVar(&actorCsvPath, "actorCsvPath", "/home/tt/.config/joj/students.csv", "path to actor csv file")
	parseMultiValueFlag(flags, &metaFile, "meta", "meta files to check")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if showVersion {
		_, err := fmt.Fprintln(stdout, Version)
		return err
	}
	setupSlog()
	slog.Info("start repo-health-checker", "version", Version)
	slog.Debug("cli args",
		"repoSize", repoSize,
		"checkFileNameList", checkFileNameList,
		"checkFileSumList", checkFileSumList,
		"whitelistedChars", whitelistedChars,
		"meta", metaFile,
	)
	res := healthcheck.All(
		rootDir,
		checkFileNameList,
		checkFileSumList,
		whitelistedChars,
		allowedDomainList,
		actorCsvPath,
		metaFile,
		repoSize,
	)
	jsonRes, err := json.Marshal(res)
	if err != nil {
		slog.Error("marshal result", "error", err)
		return err
	}
	_, err = fmt.Fprintln(stdout, string(jsonRes))
	return err
}

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		slog.Error("repo-health-checker", "error", err)
		os.Exit(1)
	}
}
