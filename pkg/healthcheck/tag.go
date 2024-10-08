package healthcheck

import (
	"fmt"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

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

func CheckTags(repoPath string, category string, n int, m int) error {
	// INFO: if category not specified, skipping this check by default
	if category == "" {
		return nil
	}
	tags, err := getTagsFromRepo(repoPath)
	if err != nil {
		return fmt.Errorf("error getting tags: %v", err)
	}
	var prefix string
	switch category {
	case "exam":
		prefix = "e"
	case "project":
		prefix = "p"
	case "homework":
		prefix = "h"
	default:
		prefix = "a"
	}
	target := prefix + fmt.Sprintf("%d", n)
	if category == "project" {
		target += fmt.Sprintf("m%d", m)
	}
	found := false
	for _, tag := range tags {
		if tag == target {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("Wrong release tag '%s' or missing release tags. Please use one of '%s'.", strings.Join(tags, "', '"), target)
	}
	return nil
}
