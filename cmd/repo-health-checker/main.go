package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/joint-online-judge/JOJ3/cmd/joj3/env"
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

var (
	rootDir           string
	repoSize          float64
	checkFileNameList string
	checkFileSumList  string
	metaFile          []string
	confPath          string
	showVersion       *bool
	Version           string
)

func init() {
	showVersion = flag.Bool("version", false, "print current version")
	flag.StringVar(&rootDir, "root", ".", "root dir for forbidden files check")
	flag.Float64Var(&repoSize, "repoSize", 2, "maximum size of the repo in MiB")
	flag.StringVar(&checkFileNameList, "checkFileNameList", "", "comma-separated list of files to check")
	flag.StringVar(&checkFileSumList, "checkFileSumList", "", "comma-separated list of expected checksums")
	parseMultiValueFlag(&metaFile, "meta", "meta files to check")
	flag.StringVar(&confPath, "confPath", "", "path to conf file for teapot check") // TODO: remove me
}

func main() {
	flag.Parse()
	if *showVersion {
		fmt.Println(Version)
		return
	}
	setupSlog()
	slog.Info("start repo-health-checker", "version", Version)
	slog.Debug("cli args",
		"repoSize", repoSize,
		"checkFileNameList", checkFileNameList,
		"checkFileSumList", checkFileSumList,
		"meta", metaFile,
	)
	groups := strings.Split(os.Getenv(env.Groups), ",")
	res := healthcheck.All(
		rootDir, checkFileNameList, checkFileSumList,
		groups, metaFile, repoSize,
	)
	jsonRes, err := json.Marshal(res)
	if err != nil {
		slog.Error("marshal result", "error", err)
		os.Exit(1)
	}
	fmt.Println(string(jsonRes))
}
