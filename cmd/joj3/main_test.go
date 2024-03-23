package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"testing"

	"focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/internal/stage"
)

func compareStageResults(t *testing.T, actual, want []stage.StageResult) {
	t.Helper()
	if len(actual) != len(want) {
		t.Fatalf("len(actual) = %d, want %d", len(actual), len(want))
	}
	for i := range actual {
		if actual[i].Name != want[i].Name {
			t.Errorf("actual[%d].Name = %s, want = %s", i, actual[i].Name,
				want[i].Name)
		}
		if len(actual[i].Results) != len(want[i].Results) {
			t.Fatalf("len(actual[%d].Results) = %d, want = %d", i,
				len(actual[i].Results), len(want[i].Results))
		}
		for j := range actual[i].Results {
			if actual[i].Results[j].Score != want[i].Results[j].Score {
				t.Errorf("actual[%d].Results[%d].Score = %d, want = %d", i, j,
					actual[i].Results[j].Score, want[i].Results[j].Score)
			}
			r := regexp.MustCompile(want[i].Results[j].Comment)
			if !r.MatchString(actual[i].Results[j].Comment) {
				t.Errorf("actual[%d].Results[%d].Comment = %s, want RegExp = %s",
					i, j, actual[i].Results[j].Comment,
					want[i].Results[j].Comment)
			}
		}
	}
}

func TestMain(t *testing.T) {
	tests := []struct {
		name string
		want []stage.StageResult
	}{
		{"success", []stage.StageResult{
			{Name: "compile", Results: []stage.ParserResult{
				{Score: 0, Comment: ""},
			}},
			{Name: "run", Results: []stage.ParserResult{
				{Score: 100, Comment: "executor status: run time: \\d+ ns, memory: \\d+ bytes"},
				{Score: 100, Comment: "executor status: run time: \\d+ ns, memory: \\d+ bytes"},
			}},
		}},
		{"compile_error", []stage.StageResult{
			{Name: "compile", Results: []stage.ParserResult{
				{Score: 0, Comment: "Unexpected executor status: Nonzero Exit Status\\."},
			}},
		}},
		{"dummy", []stage.StageResult{
			{Name: "dummy", Results: []stage.ParserResult{
				{Score: 110, Comment: "dummy comment \\+ comment from toml conf"},
			}},
		}},
		{"dummy_error", []stage.StageResult{
			{Name: "dummy", Results: []stage.ParserResult{
				{Score: 0, Comment: "Unexpected executor status: Nonzero Exit Status\\.\\s*Stderr: dummy negative score: -1"},
			}},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origDir, err := os.Getwd()
			if err != nil {
				t.Fatal(err)
			}
			err = os.Chdir(fmt.Sprintf("../../examples/%s", tt.name))
			if err != nil {
				t.Fatal(err)
			}
			defer func() {
				err := os.Chdir(origDir)
				if err != nil {
					t.Fatal(err)
				}
			}()
			os.Args = []string{"./joj3"}
			main()
			outputFile := "joj3_result.json"
			data, err := os.ReadFile(outputFile)
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(outputFile)
			var stageResults []stage.StageResult
			err = json.Unmarshal(data, &stageResults)
			if err != nil {
				t.Fatal(err)
			}
			t.Logf("stageResults: %+v", stageResults)
			compareStageResults(t, stageResults, tt.want)
		})
	}
}
