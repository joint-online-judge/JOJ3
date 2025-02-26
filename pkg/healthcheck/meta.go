package healthcheck

import (
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"strings"
)

// getMetas retrieves a list of metadata files that are expected to exist in the specified root directory.
// It checks for the existence of each file in the fileList and provides instructions if any file is missing.
func getMetas(rootDir string, fileList []string) ([]string, string, error) {
	var regexList []*regexp.Regexp
	for _, file := range fileList {
		pattern := "(?i)" + file
		if !strings.Contains(pattern, "\\.") {
			pattern += "(\\.[^\\.]*)?"
		}
		regex, err := regexp.Compile(pattern)
		if err != nil {
			return nil, "", fmt.Errorf("error compiling regex:%w", err)
		}
		regexList = append(regexList, regex)
	}
	files, err := os.ReadDir(rootDir)
	if err != nil {
		return nil, "", fmt.Errorf("error reading directory: %w", err)
	}

	matched := make([]bool, len(fileList))

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fileName := file.Name()

		for i, regex := range regexList {
			if regex.MatchString(fileName) {
				matched[i] = true
			}
		}
	}

	// Process unmatched patterns
	var unmatchedList []string
	var umatchedRes string

	for i, wasFound := range matched {
		if !wasFound {
			unmatchedList = append(unmatchedList, fileList[i])
			str := fmt.Sprintf("No %s file found", fileList[i])
			if strings.Index(strings.ToLower(fileList[i]), "readme") == 0 {
				str += ", please refer to https://www.makeareadme.com/ for more information"
			} else if strings.Index(strings.ToLower(fileList[i]), "changelog") == 0 {
				str += ", please refer to https://keepachangelog.com/en/1.1.0/ for more information"
			}
			str += ".\n"
			umatchedRes += str
		}
	}

	return unmatchedList, umatchedRes, nil
}

// metaCheck performs a check for metadata files in the specified root directory.
// It prints a message if any required metadata files are missing.
func MetaCheck(rootDir string, fileList []string) error {
	unmatchedList, umatchedRes, err := getMetas(rootDir, fileList)
	if err != nil {
		slog.Error("getting metas", "err", err)
		return fmt.Errorf("error getting metas: %w", err)
	}
	if len(unmatchedList) != 0 {
		return fmt.Errorf("%d important project file(s) missing:\n"+umatchedRes, len(unmatchedList))
	}
	return nil
}
