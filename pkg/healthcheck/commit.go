package healthcheck

import (
	"fmt"
	"log/slog"
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
