package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/joint-online-judge/JOJ3/pkg/healthcheck"
)

// parseMultiValueFlag parses a multi-value command-line flag and appends its values to the provided slice.
// It registers a flag with the specified name and description, associating it with a multiStringValue receiver.
func parseMultiValueFlag(values *[]string, flagName, description string) {
	flag.Var((*multiStringValue)(values), flagName, description)
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

// Generally, err is used for runtime errors, and checkRes is used for the result of the checks.
func main() {
	var gitWhitelist, metaFile []string
	showVersion := flag.Bool("version", false, "print current version")
	rootDir := flag.String("root", "", "")
	repo := flag.String("repo", "", "")
	localList := flag.String("localList", "", "")
	droneBranch := flag.String("droneBranch", "", "")
	releaseCategories := flag.String("releaseCategories", "", "")
	releaseNumber := flag.Int("releaseNumber", 0, "")
	checkFileNameList := flag.String("checkFileNameList", "", "Comma-separated list of files to check.")
	checkFileSumList := flag.String("checkFileSumList", "", "Comma-separated list of expected checksums.")
	parseMultiValueFlag(&gitWhitelist, "whitelist", "")
	parseMultiValueFlag(&metaFile, "meta", "")
	flag.Parse()
	if *showVersion {
		fmt.Println(Version)
		return
	}
	setupSlog()
	var err error
	err = healthcheck.RepoSize()
	if err != nil {
		fmt.Printf("## Repo Size Check Failed:\n%s\n", err.Error())
	}
	err = healthcheck.ForbiddenCheck(*rootDir, gitWhitelist, *localList, *repo, *droneBranch)
	if err != nil {
		fmt.Printf("## Forbidden File Check Failed:\n%s\n", err.Error())
	}
	err = healthcheck.MetaCheck(*rootDir, metaFile)
	if err != nil {
		fmt.Printf("## Forbidden File Check Failed:\n%s\n", err.Error())
	}
	err = healthcheck.NonAsciiFiles(*rootDir, *localList)
	if err != nil {
		fmt.Printf("## Non-ASCII Characters File Check Failed:\n%s\n", err.Error())
	}
	err = healthcheck.NonAsciiMsg(*rootDir)
	if err != nil {
		fmt.Printf("## Non-ASCII Characters Commit Message Check Failed:\n%s\n", err.Error())
	}
	err = healthcheck.CheckTags(*rootDir, *releaseCategories, *releaseNumber)
	if err != nil {
		fmt.Printf("## Release Tag Check Failed:\n%s\n", err.Error())
	}
	err = healthcheck.VerifyFiles(*rootDir, *checkFileNameList, *checkFileSumList)
	if err != nil {
		fmt.Printf("## Repo File Check Failed:\n%s\n", err.Error())
	}
}
