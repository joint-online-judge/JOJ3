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
// It searches for files that do not match the specified regex patterns in the given file list.
func getForbiddens(root string, fileList []string) ([]string, error) {
	var matches []string

	ignore, err := gitignore.NewFromFile("./.gitignore")
	if err != nil {
		return nil, err
	}

	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		} else {
			match := ignore.Relative(info.Name(), true)
			if match != nil {
				if match.Ignore() {
					matches = append(matches, path)
				}
			}
		}
		return nil
	})

	return matches, err
}

// forbiddenCheck checks for forbidden files in the specified root directory.
// It prints the list of forbidden files found, along with instructions on how to fix them.
func ForbiddenCheck(rootDir string, regexList []string) error {
	forbids, err := getForbiddens(rootDir, regexList)
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
