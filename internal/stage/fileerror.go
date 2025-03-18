package stage

import "fmt"

// copied from https://github.com/criyle/go-judge/blob/master/envexec/cmd.go
// FileErrorType defines the location that file operation fails
type FileErrorType int

// FileError enums
const (
	ErrCopyInOpenFile FileErrorType = iota
	ErrCopyInCreateDir
	ErrCopyInCreateFile
	ErrCopyInCopyContent
	ErrCopyOutOpen
	ErrCopyOutNotRegularFile
	ErrCopyOutSizeExceeded
	ErrCopyOutCreateFile
	ErrCopyOutCopyContent
	ErrCollectSizeExceeded
	ErrSymlink
)

// FileError defines the location, file name and the detailed message for a failed file operation
type FileError struct {
	Name    string        `json:"name"`
	Type    FileErrorType `json:"type"`
	Message string        `json:"message,omitempty"`
}

var fileErrorString = []string{
	"CopyInOpenFile",
	"CopyInCreateDir",
	"CopyInCreateFile",
	"CopyInCopyContent",
	"CopyOutOpen",
	"CopyOutNotRegularFile",
	"CopyOutSizeExceeded",
	"CopyOutCreateFile",
	"CopyOutCopyContent",
	"CollectSizeExceeded",
}

var fileErrorStringReverse = make(map[string]FileErrorType)

func (t FileErrorType) String() string {
	v := int(t)
	if v >= 0 && v < len(fileErrorString) {
		return fileErrorString[v]
	}
	return ""
}

// MarshalJSON encodes file error into json string
func (t FileErrorType) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.String() + `"`), nil
}

// UnmarshalJSON decodes file error from json string
func (t *FileErrorType) UnmarshalJSON(b []byte) error {
	str := string(b)
	v, ok := fileErrorStringReverse[str]
	if ok {
		return fmt.Errorf("%s is not file error type", str)
	}
	*t = v
	return nil
}

func init() {
	for i, v := range fileErrorString {
		fileErrorStringReverse[`"`+v+`"`] = FileErrorType(i)
	}
}
