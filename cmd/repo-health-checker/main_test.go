package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/joint-online-judge/JOJ3/internal/stage"
	"github.com/joint-online-judge/JOJ3/pkg/healthcheck"
)

type exampleConf struct {
	Stages []struct {
		Executor struct {
			With struct {
				Default struct {
					Args []string `json:"args"`
				} `json:"default"`
			} `json:"with"`
		} `json:"executor"`
	} `json:"stages"`
}

func readJSON[T any](t *testing.T, path string) T {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		t.Fatal(err)
	}
	return value
}

func TestHealthcheckExamples(t *testing.T) {
	root := "../../tmp/submodules/JOJ3-examples/examples/healthcheck"
	entries, err := os.ReadDir(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		// This fixture is retained in cmd/joj3 as the sandbox binary/parser smoke
		// test, so do not execute it a second time here.
		if entry.Name() == "release" {
			continue
		}
		t.Run(entry.Name(), func(t *testing.T) {
			dir := filepath.Join(root, entry.Name())
			conf := readJSON[exampleConf](t, filepath.Join(dir, "conf.json"))
			expected := readJSON[[]stage.StageResult](t, filepath.Join(dir, "expected.json"))
			if len(conf.Stages) != 1 || len(expected) != 1 || len(expected[0].Results) != 1 {
				t.Fatal("healthcheck fixture must contain one stage and one result")
			}
			args := conf.Stages[0].Executor.With.Default.Args
			if len(args) == 0 {
				t.Fatal("healthcheck fixture has no command")
			}
			t.Chdir(dir)
			var stdout bytes.Buffer
			if err := run(args[1:], &stdout); err != nil {
				t.Fatal(err)
			}
			var got healthcheck.Result
			if err := json.Unmarshal(stdout.Bytes(), &got); err != nil {
				t.Fatalf("decode output %q: %v", stdout.String(), err)
			}
			want := healthcheck.Result{
				Msg:    expected[0].Results[0].Comment,
				Failed: expected[0].ForceQuit,
			}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("run() = %+v, want %+v", got, want)
			}
		})
	}
}

func TestRunVersionAndInvalidFlag(t *testing.T) {
	oldVersion := Version
	Version = "test-version"
	t.Cleanup(func() { Version = oldVersion })
	var stdout bytes.Buffer
	if err := run([]string{"-version"}, &stdout); err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(stdout.String()) != Version {
		t.Fatalf("version output = %q", stdout.String())
	}
	if err := run([]string{"-unknown"}, &stdout); err == nil {
		t.Fatal("run() accepted an unknown flag")
	}
}
