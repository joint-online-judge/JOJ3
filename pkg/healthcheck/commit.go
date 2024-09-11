package healthcheck

import (
	"fmt"
	"log/slog"
	"strings"
	"unicode"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
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
	commits, err := repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		slog.Error("getting commits", "err", err)
		return fmt.Errorf("error getting commits from reference %s: %v", ref.Hash(), err)
	}

	var msgs []string
	err = commits.ForEach(func(c *object.Commit) error {
		msgs = append(msgs, c.Message)
		return nil
	})
	if err != nil {
		slog.Error("iterating commits", "err", err)
		return fmt.Errorf("error iterating commits: %v", err)
	}

	var nonAsciiMsgs []string
	for _, msg := range msgs {
		if msg == "" {
			continue
		}
		if strings.HasPrefix(msg, "Merge pull request") {
			continue
		}
		for _, c := range msg {
			if c > unicode.MaxASCII {
				nonAsciiMsgs = append(nonAsciiMsgs, msg)
			}
		}
	}
	if len(nonAsciiMsgs) > 0 {
		return fmt.Errorf("Non-ASCII characters in commit messages:\n%s", strings.Join(nonAsciiMsgs, "\n"))
	}
	return nil
}
