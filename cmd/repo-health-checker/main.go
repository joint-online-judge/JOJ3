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
	var metaFile []string
	showVersion := flag.Bool("version", false, "print current version")
	rootDir := flag.String("root", ".", "root dir for forbidden files check")
	repoSize := flag.Float64("repoSize", 2, "maximum size of the repo in MiB")
	// TODO: remove localList, it is only for backward compatibility now
	localList := flag.String("localList", "", "local file list for non-ascii file check")
	checkFileNameList := flag.String("checkFileNameList", "", "comma-separated list of files to check")
	checkFileSumList := flag.String("checkFileSumList", "", "comma-separated list of expected checksums")
	parseMultiValueFlag(&metaFile, "meta", "meta files to check")
	// TODO: remove gitWhitelist, it is only for backward compatibility now
	var gitWhitelist []string
	parseMultiValueFlag(&gitWhitelist, "whitelist", "[DEPRECATED] will be ignored")
	flag.Parse()
	if *showVersion {
		fmt.Println(Version)
		return
	}
	setupSlog()
	slog.Info("start repo-health-checker", "version", Version)
	slog.Debug("cli args",
		"repoSize", repoSize,
		"localList", localList,
		"checkFileNameList", checkFileNameList,
		"checkFileSumList", checkFileSumList,
		"meta", metaFile,
	)
	var err error
	err = healthcheck.RepoSize(*repoSize)
	if err != nil {
		fmt.Printf("### Repo Size Check Failed:\n%s\n", err.Error())
	}
	err = healthcheck.ForbiddenCheck(*rootDir)
	if err != nil {
		fmt.Printf("### Forbidden File Check Failed:\n%s\n", err.Error())
	}
	err = healthcheck.MetaCheck(*rootDir, metaFile)
	if err != nil {
		fmt.Printf("### Forbidden File Check Failed:\n%s\n", err.Error())
	}
	err = healthcheck.NonAsciiFiles(*rootDir)
	if err != nil {
		fmt.Printf("### Non-ASCII Characters File Check Failed:\n%s\n", err.Error())
	}
	err = healthcheck.NonAsciiMsg(*rootDir)
	if err != nil {
		fmt.Printf("### Non-ASCII Characters Commit Message Check Failed:\n%s\n", err.Error())
	}
	err = healthcheck.VerifyFiles(*rootDir, *checkFileNameList, *checkFileSumList)
	if err != nil {
		fmt.Printf("### Repo File Check Failed:\n%s\n", err.Error())
	}
}
