package healthcheck

import (
	"fmt"
	"log/slog"
	"strings"
	"unicode"

	"github.com/go-git/go-git/v5"
)

// nonAsciiMsg checks for non-ASCII characters in the commit message.
// If the message starts with "Merge pull request", it skips the non-ASCII characters check.
// Otherwise, it iterates over each character in the message and checks if it is a non-ASCII character.
// If a non-ASCII character is found, it returns an error indicating not to use non-ASCII characters in commit messages.
// Otherwise, it returns nil indicating that the commit message is valid.
func NonAsciiMsg(root string) error {
	// cmd := exec.Command("git", "log", "--encoding=UTF-8", "--format=%B")
	repo, err := git.PlainOpen(root)
	if err != nil {
		slog.Error("openning git repo", "err", err)
		return fmt.Errorf("error openning git repo: %v", err)
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
	if msg == "" {
		return nil
	}

	var isCommitLegal bool = true
	// List of prefixes to ignore in the commit message
	ignoredPrefixes := []string{
		"Co-authored-by:",
		"Reviewed-by:",
		"Co-committed-by:",
		"Reviewed-on:",
	}

	// Split message by lines and ignore specific lines with prefixes
	lines := strings.Split(msg, "\n")
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		ignore := false
		for _, prefix := range ignoredPrefixes {
			if strings.HasPrefix(trimmedLine, prefix) {
				ignore = true
				break
			}
		}
		if ignore {
			continue
		}
		// Check for non-ASCII characters in the rest of the lines
		for _, c := range line {
			if c > unicode.MaxASCII {
				isCommitLegal = false
				break
			}
		}
	}

	if !isCommitLegal {
		return fmt.Errorf("Non-ASCII characters in commit messages:\n%s", msg)
	}
	return nil
}
