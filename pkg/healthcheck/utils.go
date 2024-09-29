package healthcheck

import (
	"fmt"
	"regexp"
)

func inString(str1 string, strList []string) bool {
	for _, str := range strList {
		if str1 == str {
			return true
		}
	}
	return false
}

// addExt appends the specified extension to each file name in the given fileList.
// It modifies the original fileList in place.
func addExt(fileList []string, ext string) {
	for i, file := range fileList {
		fileList[i] = file + ext
	}
}

// getRegex compiles each regex pattern in the fileList into a []*regexp.Regexp slice.
// It returns a slice containing compiled regular expressions.
func getRegex(fileList []string) ([]*regexp.Regexp, error) {
	var regexList []*regexp.Regexp
	for _, pattern := range fileList {
		regex, err := regexp.Compile("(?i)" + pattern)
		if err != nil {
			return nil, fmt.Errorf("error compiling regex: %w", err)
		}
		regexList = append(regexList, regex)
	}

	return regexList, nil
}
