package healthcheck

import (
	"encoding/csv"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"unicode"

	"github.com/go-git/go-git/v5"
)

func checkMsg(msg string) bool {
	// List of prefixes to ignore in the commit message
	ignoredPrefixes := []string{
		"Co-authored-by:",
		"Reviewed-by:",
		"Co-committed-by:",
		"Reviewed-on:",
	}

	// Split message by lines and ignore specific lines with prefixes
	lines := strings.Split(msg, "\n")
	for i, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		ignore := false
		if i != 0 {
			for _, prefix := range ignoredPrefixes {
				if strings.HasPrefix(trimmedLine, prefix) {
					ignore = true
					break
				}
			}
		}
		if ignore {
			continue
		}
		// Check for non-ASCII characters in the rest of the lines
		for _, c := range line {
			if c > unicode.MaxASCII {
				return false
			}
		}
	}
	return true
}

// NonASCIIMsg checks for non-ASCII characters in the commit message.
// It iterates over each character in the message and checks if it is a non-ASCII character.
// If a non-ASCII character is found, it returns an error indicating not to use non-ASCII characters in commit messages.
// Otherwise, it returns nil indicating that the commit message is valid.
// It skips the non-ASCII characters check for lines starting with specific keywords like "Co-authored-by", "Reviewed-by", and "Co-committed-by".
func NonASCIIMsg(root string) error {
	repo, err := git.PlainOpen(root)
	if err != nil {
		slog.Error("opening git repo", "err", err)
		return fmt.Errorf("error opening git repo: %v", err)
	}

	ref, err := repo.Head()
	if err != nil {
		slog.Error("getting reference", "err", err)
		return fmt.Errorf("error getting reference: %v", err)
	}

	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		slog.Error("getting latest commit", "err", err)
		return fmt.Errorf("error getting latest commit: %v", err)
	}

	msg := commit.Message
	if !checkMsg(msg) {
		return fmt.Errorf("Non-ASCII characters in commit messages:\n%s", msg)
	}
	return nil
}

func AuthorEmailCheck(root string, allowedDomains []string, actorCsvPath string) error {
	repo, err := git.PlainOpen(root)
	if err != nil {
		slog.Error("opening git repo", "err", err)
		return fmt.Errorf("error opening git repo: %v", err)
	}

	ref, err := repo.Head()
	if err != nil {
		slog.Error("getting reference", "err", err)
		return fmt.Errorf("error getting reference: %v", err)
	}

	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		slog.Error("getting latest commit", "err", err)
		return fmt.Errorf("error getting latest commit: %v", err)
	}

	email := strings.ToLower(commit.Author.Email)

	checkDomains := true
	if _, err := os.Stat(actorCsvPath); err != nil {
		slog.Error("checking actor CSV file stat", "err", err, "path", actorCsvPath)
	} else {
		f, err := os.Open(actorCsvPath)
		if err != nil {
			slog.Error("opening actor CSV file", "err", err, "path", actorCsvPath)
		} else {
			defer f.Close()
			reader := csv.NewReader(f)
			for {
				row, err := reader.Read()
				if err == io.EOF {
					checkDomains = false
					break
				}
				if err != nil {
					slog.Error("reading actor CSV file", "err", err, "path", actorCsvPath)
					break
				}
				if len(row) >= 3 {
					actor := row[2]
					for _, domain := range allowedDomains {
						if email == actor+"@"+domain {
							return nil
						}
					}
				}
			}
		}
	}
	var msgPrefix string
	if checkDomains {
		for _, domain := range allowedDomains {
			if strings.HasSuffix(email, "@"+domain) {
				return nil
			}
		}
		msgPrefix = fmt.Sprintf(
			"Author email %s is not in the allowed domains: `%s`",
			email,
			strings.Join(allowedDomains, "`, `"),
		)
	} else {
		msgPrefix = fmt.Sprintf(
			"Author email %s is not stored in `%s`",
			email,
			actorCsvPath,
		)
	}

	return fmt.Errorf("%s\n\n"+
		"To fix it, please run the following commands:\n"+
		"```bash\n"+
		"git config user.email \"<your_email>\"\n"+
		"git commit --amend --no-edit --reset-author\n"+
		"git push --force\n"+
		"```\n"+
		"Replace `<your_email>` with your actual email address\n",
		msgPrefix,
	)
}
