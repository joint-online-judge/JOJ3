package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/joint-online-judge/JOJ3/internal/stage"
)

func compareStageResults(t *testing.T, actual, expected []stage.StageResult, regex bool) {
	t.Helper()
	if len(actual) != len(expected) {
		t.Fatalf("len(actual) = %d, expected %d", len(actual), len(expected))
	}
	for i := range actual {
		if actual[i].Name != expected[i].Name {
			t.Errorf("actual[%d].Name = %s, expected = %s", i, actual[i].Name,
				expected[i].Name)
		}
		if len(actual[i].Results) != len(expected[i].Results) {
			t.Fatalf("len(actual[%d].Results) = %d, expected = %d", i,
				len(actual[i].Results), len(expected[i].Results))
		}
		for j := range actual[i].Results {
			if actual[i].Results[j].Score != expected[i].Results[j].Score {
				t.Errorf("actual[%d].Results[%d].Score = %d, expected = %d", i, j,
					actual[i].Results[j].Score, expected[i].Results[j].Score)
			}
			if regex {
				r := regexp.MustCompile(expected[i].Results[j].Comment)
				if !r.MatchString(actual[i].Results[j].Comment) {
					t.Errorf("actual[%d].Results[%d].Comment = %s, expected RegExp = %s",
						i, j, actual[i].Results[j].Comment,
						expected[i].Results[j].Comment)
				}
			} else if actual[i].Results[j].Comment != expected[i].Results[j].Comment {
				t.Errorf("actual[%d].Results[%d].Comment = %s, expected = %s", i, j,
					actual[i].Results[j].Comment, expected[i].Results[j].Comment)
			}
		}
	}
}

func readStageResults(t *testing.T, path string) []stage.StageResult {
	t.Helper()
	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	var results []stage.StageResult
	err = json.NewDecoder(file).Decode(&results)
	if err != nil {
		t.Fatal(err)
	}
	return results
}

func TestRun(t *testing.T) {
	var tests []string
	root := "../../tmp/submodules/JOJ3-examples/examples/"
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if path == root {
				return nil
			}
			path0 := filepath.Join(path, "expected_regex.json")
			path1 := filepath.Join(path, "expected.json")
			_, err0 := os.Stat(path0)
			_, err1 := os.Stat(path1)
			if err0 != nil && err1 != nil {
				return nil
			}
			tests = append(tests, strings.TrimPrefix(path, root))
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range tests {
		// The repo-health-checker package runs all healthcheck fixtures directly so
		// their execution contributes to Go coverage. Keep one sandbox case here as
		// an end-to-end binary/executor/parser smoke test.
		if strings.HasPrefix(tt, "healthcheck/") && tt != "healthcheck/release" {
			continue
		}
		t.Run(tt, func(t *testing.T) {
			if tt == "healthcheck/release" {
				prepareLargeCopyInFixture(t, filepath.Join(root, tt))
			}
			t.Chdir(fmt.Sprintf("%s%s", root, tt))
			os.Args = []string{"./joj3"}
			outputFile := "joj3_result.json"
			defer os.Remove(outputFile)
			runningTest = true
			_ = mainImpl()
			stageResults := readStageResults(t, outputFile)
			regex := true
			expectedFile := "expected_regex.json"
			if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
				regex = false
				expectedFile = "expected.json"
			}
			expectedStageResults := readStageResults(t, expectedFile)
			compareStageResults(t, stageResults, expectedStageResults, regex)
		})
	}
}

func prepareLargeCopyInFixture(t *testing.T, dir string) {
	t.Helper()
	const fileCount = 1001
	filesDir := filepath.Join(dir, "many-files")
	if err := os.Mkdir(filesDir, 0o700); err != nil && !os.IsExist(err) {
		t.Fatal(err)
	}
	for i := range fileCount {
		path := filepath.Join(filesDir, fmt.Sprintf("%04d", i))
		if err := os.WriteFile(path, nil, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	t.Cleanup(func() { _ = os.RemoveAll(filesDir) })

	confPath := filepath.Join(dir, "conf.json")
	original, err := os.ReadFile(confPath)
	if err != nil {
		t.Fatal(err)
	}
	var conf map[string]any
	if err := json.Unmarshal(original, &conf); err != nil {
		t.Fatal(err)
	}
	stages := conf["stages"].([]any)
	executor := stages[0].(map[string]any)["executor"].(map[string]any)
	with := executor["with"].(map[string]any)
	command := with["default"].(map[string]any)
	args := command["args"].([]any)
	checkerArgs := make([]string, 0, len(args))
	for _, arg := range args {
		checkerArgs = append(checkerArgs, fmt.Sprintf("%q", arg))
	}
	command["args"] = []string{
		"/bin/sh", "-c",
		fmt.Sprintf("test \"$(find many-files -type f | wc -l)\" -eq %d && exec %s", fileCount, strings.Join(checkerArgs, " ")),
	}
	patched, err := json.Marshal(conf)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(confPath, patched, 0o600); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.WriteFile(confPath, original, 0o600); err != nil {
			t.Errorf("restore conf: %v", err)
		}
	})
}
