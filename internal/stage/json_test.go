package stage

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestStatusJSONRoundTrip(t *testing.T) {
	for status := StatusInvalid; status <= StatusInternalError; status++ {
		data, err := json.Marshal(status)
		if err != nil {
			t.Fatal(err)
		}
		var got Status
		if err := json.Unmarshal(data, &got); err != nil {
			t.Fatalf("unmarshal %s: %v", data, err)
		}
		if got != status {
			t.Fatalf("round trip = %v, want %v", got, status)
		}
	}
	var status Status
	if err := json.Unmarshal([]byte(`"unknown"`), &status); err == nil {
		t.Fatal("unknown status was accepted")
	}
}

func TestFileErrorTypeJSONRoundTrip(t *testing.T) {
	for fileError := ErrCopyInOpenFile; fileError <= ErrCollectSizeExceeded; fileError++ {
		data, err := json.Marshal(fileError)
		if err != nil {
			t.Fatal(err)
		}
		var got FileErrorType
		if err := json.Unmarshal(data, &got); err != nil {
			t.Fatalf("unmarshal %s: %v", data, err)
		}
		if got != fileError {
			t.Fatalf("round trip = %v, want %v", got, fileError)
		}
	}
	var fileError FileErrorType
	if err := json.Unmarshal([]byte(`"unknown"`), &fileError); err == nil {
		t.Fatal("unknown file error type was accepted")
	}
}

func TestExecutorResultJSONSummarizesFileContents(t *testing.T) {
	data, err := json.Marshal(ExecutorResult{
		Status: StatusAccepted,
		Files:  map[string]string{"stdout": "secret output"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), "secret output") || !strings.Contains(string(data), `"stdout":"len:13"`) {
		t.Fatalf("MarshalJSON() = %s", data)
	}
}

func TestNonNullSliceMarshalsEmptyAsArray(t *testing.T) {
	data, err := json.Marshal(NonNullSlice[int](nil))
	if err != nil || string(data) != "[]" {
		t.Fatalf("MarshalJSON() = %s, %v", data, err)
	}
}
