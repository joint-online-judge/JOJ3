package healthcheck

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/format/gitattributes"
)

// getSubmodulePathsFromGoGit uses the go-git library to open the repository
// at the given root path and retrieve a list of all submodule paths.
// It returns a set of submodule paths for efficient lookup.
func getSubmodulePathsFromGoGit(root string) (map[string]struct{}, error) {
	submodulePaths := make(map[string]struct{})

	// Open the git repository at the given path.
	repo, err := git.PlainOpen(root)
	if err != nil {
		if err == git.ErrRepositoryNotExists {
			return submodulePaths, nil
		}
		return nil, fmt.Errorf("error opening git repository: %w", err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return nil, fmt.Errorf("error getting worktree: %w", err)
	}

	// Get the list of submodules.
	submodules, err := worktree.Submodules()
	if err != nil {
		return nil, fmt.Errorf("error getting submodules: %w", err)
	}

	for _, sm := range submodules {
		submodulePaths[filepath.ToSlash(sm.Config().Path)] = struct{}{}
	}

	return submodulePaths, nil
}

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

	submodules, err := getSubmodulePathsFromGoGit(root)
	if err != nil {
		return nil, err
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

		relPath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		if info.IsDir() {
			if info.Name() == ".git" {
				return filepath.SkipDir
			}
			if _, isSubmodule := submodules[relPath]; isSubmodule {
				return filepath.SkipDir
			}
			return nil
		}

		if gitattrExist {
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
