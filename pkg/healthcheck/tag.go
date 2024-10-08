package healthcheck

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/joint-online-judge/JOJ3/cmd/joj3/conf"
)

func parseConventionalCommit(commit string) (*conf.ConventionalCommit, error) {
	re := regexp.MustCompile(`(?s)^(\w+)(\(([^)]+)\))?!?: (.+?)(\n\n(.+?))?(\n\n(.+))?$`)
	matches := re.FindStringSubmatch(strings.TrimSpace(commit))
	if matches == nil {
		return nil, fmt.Errorf("invalid conventional commit format")
	}
	cc := &conf.ConventionalCommit{
		Type:        matches[1],
		Scope:       matches[3],
		Description: strings.TrimSpace(matches[4]),
		Body:        strings.TrimSpace(matches[6]),
		Footer:      strings.TrimSpace(matches[8]),
	}
	return cc, nil
}

func getTagFromMsg() (tag string, err error) {
	msg, err := conf.GetCommitMsg()
	if err != nil {
		return "", err
	}
	conventionalCommit, err := parseConventionalCommit(msg)
	if err != nil {
		return "", err
	}
	return conventionalCommit.Scope, err
}

func getTagsFromRepo(repoPath string) ([]string, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, fmt.Errorf("error opening repo: %v", err)
	}

	refs, err := repo.Tags()
	if err != nil {
		return nil, fmt.Errorf("error getting tags: %v", err)
	}

	var tags []string
	err = refs.ForEach(func(ref *plumbing.Reference) error {
		tags = append(tags, ref.Name().Short())
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error iterating tags: %v", err)
	}

	return tags, nil
}

func CheckTags(repoPath string, skip bool) error {
	if skip {
		return nil
	}
	tags, err := getTagsFromRepo(repoPath)
	if err != nil {
		return fmt.Errorf("error getting tags from repo: %v", err)
	}

	target, err := getTagFromMsg()
	if err != nil {
		return fmt.Errorf("error getting tag from msg scope: %v", err)
	}
	found := false
	for _, tag := range tags {
		if tag == target {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("Wrong release tag in '%s' or missing release tags. Please use '%s'.", strings.Join(tags, "', '"), target)
	}
	return nil
}
