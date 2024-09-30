package healthcheck

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

// getNonAscii retrieves a list of files in the specified root directory that contain non-ASCII characters.
// It searches for non-ASCII characters in each file's content and returns a list of paths to files containing non-ASCII characters.
func getNonAscii(root string, localList string) ([]string, error) {
	var nonAscii []string

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

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if info.Name() == ".git" || info.Name() == ".gitea" || info.Name() == "ci" || (localList != "" && inString(info.Name(), dirs)) {
				return filepath.SkipDir
			} else {
				return nil
			}
		}

		if info.Name() == "healthcheck" {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			cont := true
			for _, c := range scanner.Text() {
				if c > unicode.MaxASCII {
					nonAscii = append(nonAscii, "\t"+path)
					cont = false
					break
				}
			}
			if !cont {
				break
			}
		}

		return nil
	})

	return nonAscii, err
}

// nonAsciiFiles checks for non-ASCII characters in files within the specified root directory.
// It prints a message with the paths to files containing non-ASCII characters, if any.
func NonAsciiFiles(root string, localList string) error {
	nonAscii, err := getNonAscii(root, localList)
	if err != nil {
		slog.Error("getting non-ascii", "err", err)
		return fmt.Errorf("error getting non-ascii: %w", err)
	}
	if len(nonAscii) > 0 {
		return fmt.Errorf("Non-ASCII characters found in the following files:\n%s",
			strings.Join(nonAscii, "\n"))
	}
	return nil
}
