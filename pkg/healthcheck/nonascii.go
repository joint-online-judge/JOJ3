package healthcheck

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/go-git/go-git/v5/plumbing/format/gitattributes"
)

// getNonASCII retrieves a list of files in the specified root directory that contain non-ASCII characters.
// It searches for non-ASCII characters in each file's content and returns a list of paths to files containing non-ASCII characters.
func getNonASCII(root string) ([]string, error) {
	var nonASCII []string
	gitattrExist := true
	var matcher gitattributes.Matcher
	_, err := os.Stat(".gitattributes")
	if os.IsNotExist(err) {
		gitattrExist = false
	}

	if gitattrExist {
		fs := os.DirFS(".")
		f, err := fs.Open(".gitattributes")
		if err != nil {
			return nil, err
		}

		attribute, err := gitattributes.ReadAttributes(f, nil, true)
		if err != nil {
			return nil, err
		}
		matcher = gitattributes.NewMatcher(attribute)
	}

	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if info.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}

		if gitattrExist {
			relPath, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			ret, matched := matcher.Match(strings.Split(relPath, "/"), nil)
			if matched && ret["text"].IsUnset() && !ret["text"].IsSet() {
				return nil
			}
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
					nonASCII = append(nonASCII, "\t"+path)
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

	return nonASCII, err
}

// NonASCIIFiles checks for non-ASCII characters in files within the specified root directory.
// It prints a message with the paths to files containing non-ASCII characters, if any.
func NonASCIIFiles(root string) error {
	nonASCII, err := getNonASCII(root)
	if err != nil {
		slog.Error("getting non-ascii", "err", err)
		return fmt.Errorf("error getting non-ascii: %w", err)
	}
	if len(nonASCII) > 0 {
		return fmt.Errorf("Non-ASCII characters found in the following files:\n%s",
			strings.Join(nonASCII, "\n"))
	}
	return nil
}
