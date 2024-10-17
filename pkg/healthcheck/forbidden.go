package healthcheck

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/denormal/go-gitignore"
)

// getForbiddens retrieves a list of forbidden files in the specified root directory.
// It searches for files that match the specified ignore patterns in the .gitignore file.
func getForbiddens(root string) ([]string, error) {
	var matches []string

	// Create a gitignore instance from the .gitignore file
	ignore := gitignore.NewRepositoryWithCache(root, ".gitignore", gitignore.NewCache(), func(e gitignore.Error) bool {
		return false
	})

	var err error

	if err != nil {
		return nil, err
	}

	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if info.Name() == ".git" {
				return filepath.SkipDir
			} else if info.Name() == root {
				return nil
			}
		}

		// Get the relative path to the git repo root
		relPath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		match := ignore.Relative(relPath, true)

		// Check if the relative file path should be ignored based on the .gitignore rules
		if match != nil && match.Ignore() {
			matches = append(matches, path)
		}

		return nil
	})

	return matches, err
}

// ForbiddenCheck checks for forbidden files in the specified root directory.
// It prints the list of forbidden files found, along with instructions on how to fix them.
func ForbiddenCheck(rootDir string) error {
	forbids, err := getForbiddens(rootDir)
	if err != nil {
		slog.Error("getting forbiddens", "error", err)
		return fmt.Errorf("error getting forbiddens: %w", err)
	}

	if len(forbids) > 0 {
		return fmt.Errorf("The following forbidden files were found: `%s`\n\n"+
			"To fix it, first make a backup of your repository and then run the following commands:\n"+
			"```bash\n"+
			"export GIT_BRANCH=$(git branch --show-current)\n"+
			"export GIT_REMOTE_URL=$(git config --get remote.origin.url)\n"+
			"for i in %s; do git filter-repo --force --invert-paths --path \"$i\"; done\n"+
			"git remote add origin $GIT_REMOTE_URL\n"+
			"git push --set-upstream origin $GIT_BRANCH --force\n"+
			"```\n",
			strings.Join(forbids, "`, `"),
			strings.Join(forbids, " "))
	}
	return nil
}
