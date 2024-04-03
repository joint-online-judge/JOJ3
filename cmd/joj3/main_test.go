package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"testing"

	"focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/internal/stage"
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

func TestMain(t *testing.T) {
	tests := []string{
		"compile/success",
		"compile/error",
		"dummy/success",
		"dummy/error",
		"cpplint/sillycode",
	}
	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			origDir, err := os.Getwd()
			if err != nil {
				t.Fatal(err)
			}
			err = os.Chdir(fmt.Sprintf("../../examples/%s", tt))
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
			outputFile := "joj3_result.json"
			defer os.Remove(outputFile)
			main()
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
