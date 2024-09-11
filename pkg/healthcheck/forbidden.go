package healthcheck

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// getForbiddens retrieves a list of forbidden files in the specified root directory.
// It searches for files that do not match the specified regex patterns in the given file list.
func getForbiddens(root string, fileList []string, localList string) ([]string, error) {
	var matches []string

	var regexList []*regexp.Regexp
	regexList, err := getRegex(fileList)
	if err != nil {
		return nil, err
	}

	var dirs []string

	if localList != "" {
		file, err := os.Open(localList)
		if err != nil {
			return nil, fmt.Errorf("Failed to open file %s: %v\n", localList, err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			dirs = append(dirs, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("Error reading file %s: %v\n", localList, err)
		}
	}

	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if info.Name() == ".git" || info.Name() == ".gitea" || info.Name() == "ci" || (localList != "" && inString(info.Name(), dirs)) {
				return filepath.SkipDir
			}
		} else {
			match := false
			for _, regex := range regexList {
				if regex.MatchString(info.Name()) {
					match = true
					break
				}
			}

			if !match {
				matches = append(matches, path)
			}
		}

		return nil
	})

	return matches, err
}

// forbiddenCheck checks for forbidden files in the specified root directory.
// It prints the list of forbidden files found, along with instructions on how to fix them.
func ForbiddenCheck(rootDir string, regexList []string, localList string, repo string, droneBranch string) error {
	forbids, err := getForbiddens(rootDir, regexList, localList)
	if err != nil {
		slog.Error("getting forbiddens", "error", err)
		return fmt.Errorf("error getting forbiddens: %w", err)
	}

	if len(forbids) > 0 {
		return fmt.Errorf("The following forbidden files were found: %s\n\nTo fix it, first make a backup of your repository and then run the following commands:\nfor i in %s%s",
			strings.Join(forbids, ", "),
			strings.Join(forbids, " "),
			fmt.Sprint(
				"; do git filter-repo --force --invert-paths --path \"$i\"; done\ngit remote add origin ",
				repo, "\ngit push --set-upstream origin ",
				droneBranch, " --force"))
	}
	return nil
}
