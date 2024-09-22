package healthcheck

import (
	"fmt"
	"log/slog"
	"os"
)

// getMetas retrieves a list of metadata files that are expected to exist in the specified root directory.
// It checks for the existence of each file in the fileList and provides instructions if any file is missing.
func getMetas(rootDir string, fileList []string) ([]string, string, error) {
	addExt(fileList, "\\.*")
	regexList, err := getRegex(fileList)
	var unmatchedList []string

	if err != nil {
		return nil, "", err
	}

	files, err := os.ReadDir(rootDir)
	if err != nil {
		return nil, "", fmt.Errorf("error reading directory: %w", err)
	}

	matched := false
	umatchedRes := ""

	// TODO: it seems that there is no good find substitution now
	// modify current code if exist a better solution
	for i, regex := range regexList {
		for _, file := range files {
			if file.IsDir() {
				continue
			}

			if regex.MatchString(file.Name()) {
				matched = true
				break
			}
		}
		if !matched {
			unmatchedList = append(unmatchedList, fileList[i])
			str := fmt.Sprint("\tno ", fileList[i], " file found")
			switch fileList[i] {
			case "readme\\.*":
				str += ", please refer to https://www.makeareadme.com/ for more information"
			case "changelog\\.*":
				str += ", please refer to https://keepachangelog.com/en/1.1.0/ for more information"
			default:
				str += ""
			}
			str += "\n"

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
		return fmt.Errorf("%d important project files missing\n"+umatchedRes, len(unmatchedList))
	}
	return nil
}
