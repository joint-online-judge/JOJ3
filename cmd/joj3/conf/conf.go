// Package conf provides a configuration file parser for JOJ3.
// The configuration file path is determined by the commit message.
package conf

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/koding/multiconfig"
)

func GetCommitMsg() (msg string, err error) {
	r, err := git.PlainOpen(".")
	if err != nil {
		return
	}
	ref, err := r.Head()
	if err != nil {
		return
	}
	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		return
	}
	msg = commit.Message
	return
}

func parseConventionalCommit(commit string) (*ConventionalCommit, error) {
	re := regexp.MustCompile(`(?s)^(\w+)(\(([^)]+)\))?!?: (.+?(\[([^\]]+)\])?)(\n\n(.+?))?(\n\n(.+))?$`)
	matches := re.FindStringSubmatch(strings.TrimSpace(commit))
	if matches == nil {
		return &ConventionalCommit{}, fmt.Errorf("invalid conventional commit format")
	}
	cc := &ConventionalCommit{
		Type:        matches[1],
		Scope:       matches[3],
		Description: strings.TrimSpace(matches[4]),
		Group:       matches[6],
		Body:        strings.TrimSpace(matches[8]),
		Footer:      strings.TrimSpace(matches[10]),
	}
	return cc, nil
}

func ParseConfFile(path string) (conf *Conf, err error) {
	conf = new(Conf)
	d := &multiconfig.DefaultLoader{}
	d.Loader = multiconfig.MultiLoader(
		&multiconfig.TagLoader{},
		&multiconfig.JSONLoader{Path: path},
	)
	d.Validator = multiconfig.MultiValidator(&multiconfig.RequiredValidator{})
	if err = d.Load(conf); err != nil {
		slog.Error("parse stages conf", "error", err)
		return
	}
	if err = d.Validate(conf); err != nil {
		slog.Error("validate stages conf", "error", err)
		return
	}
	return
}

func GetSHA256(filePath string) (hashStr string, err error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer func() {
		if cerr := file.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	// Calculate SHA-256
	hash := sha256.New()
	if _, err = io.Copy(hash, file); err != nil {
		return
	}
	hashStr = hex.EncodeToString(hash.Sum(nil))
	return hashStr, nil
}

func parseMsg(confRoot, confName, msg, tag string) (
	confPath string, conventionalCommit *ConventionalCommit, err error,
) {
	slog.Info("parse msg", "msg", msg)
	if tag == "" {
		conventionalCommit, err = parseConventionalCommit(msg)
		if err != nil {
			return
		}
	} else {
		conventionalCommit = &ConventionalCommit{
			Scope: tag,
			Group: "all",
		}
	}
	slog.Info("conventional commit", "commit", conventionalCommit)
	confRoot = filepath.Clean(confRoot)
	confPath = filepath.Join(confRoot, conventionalCommit.Scope, confName)
	relPath, err := filepath.Rel(confRoot, confPath)
	if err != nil {
		return
	}
	if strings.HasPrefix(relPath, "..") || filepath.IsAbs(relPath) {
		err = fmt.Errorf("invalid scope as path: %s", conventionalCommit.Scope)
		return
	}
	return
}

func hintValidScopes(confRoot, confName string) {
	confRoot = filepath.Clean(confRoot)
	validScopes := []string{}
	_ = filepath.WalkDir(confRoot, func(
		path string, d fs.DirEntry, err error,
	) error {
		if err != nil {
			slog.Error("list valid scopes", "error", err)
			return err
		}
		if d.IsDir() {
			confPath := filepath.Join(path, confName)
			if _, err := os.Stat(confPath); err == nil {
				relPath, err := filepath.Rel(confRoot, path)
				if err != nil {
					return err
				}
				if relPath == "." {
					relPath = ""
				}
				validScopes = append(validScopes,
					fmt.Sprintf("'%s'", relPath))
			}
		}
		return nil
	})
	slog.Info("HINT: use valid scopes in commit message",
		"valid scopes", validScopes)
}

func GetConfPath(confRoot, confName, fallbackConfName, msg, tag string) (
	confPath string, confStat fs.FileInfo,
	conventionalCommit *ConventionalCommit, err error,
) {
	confPath, conventionalCommit, err = parseMsg(confRoot, confName, msg, tag)
	if err != nil {
		slog.Error("parse msg", "error", err)
		// no fallback when tag is specified, it is triggered by release.yaml
		if os.IsNotExist(err) {
			slog.Info("tag is not empty, no fallback conf")
			hintValidScopes(confRoot, confName)
			return confPath, confStat, conventionalCommit, err
		}
		// fallback to conf file in conf root on parse error
		confPath = filepath.Join(confRoot, fallbackConfName)
		slog.Info("fallback to conf", "path", confPath)
	}
	confStat, err = os.Stat(confPath)
	if err != nil {
		if os.IsNotExist(err) {
			hintValidScopes(confRoot, confName)
		}
		slog.Error("stat conf", "error", err)
		// fallback to conf file in conf root on conf not exist
		confPath = filepath.Join(confRoot, fallbackConfName)
		slog.Info("fallback to conf", "path", confPath)
		confStat, err = os.Stat(confPath)
		if err != nil {
			slog.Error("stat fallback conf", "error", err)
			return confPath, confStat, conventionalCommit, err
		}
	}
	return confPath, confStat, conventionalCommit, err
}

func MatchGroups(conf *Conf, conventionalCommit *ConventionalCommit) []string {
	seen := make(map[string]bool)
	keywords := []string{}
	loweredCommitGroup := strings.ToLower(conventionalCommit.Group)
	if loweredCommitGroup == "all" {
		for i := range conf.PreStages {
			conf.PreStages[i].Group = ""
			conf.PreStages[i].Groups = nil
		}
		for i := range conf.Stages {
			conf.Stages[i].Group = ""
			conf.Stages[i].Groups = nil
		}
		for i := range conf.PostStages {
			conf.PostStages[i].Group = ""
			conf.PostStages[i].Groups = nil
		}
	}
	confStages := []ConfStage{}
	confStages = append(confStages, conf.PreStages...)
	confStages = append(confStages, conf.Stages...)
	confStages = append(confStages, conf.PostStages...)
	for _, stage := range confStages {
		if stage.Group != "" {
			keyword := strings.ToLower(stage.Group)
			if _, exists := seen[keyword]; !exists {
				seen[keyword] = true
				keywords = append(keywords, keyword)
			}
		}
		if len(stage.Groups) > 0 {
			for _, group := range stage.Groups {
				keyword := strings.ToLower(group)
				if _, exists := seen[keyword]; !exists {
					seen[keyword] = true
					keywords = append(keywords, keyword)
				}
			}
		}
	}
	slog.Info("group keywords from stages", "keywords", keywords)
	groups := []string{}
	for _, keyword := range keywords {
		if strings.Contains(loweredCommitGroup, keyword) {
			groups = append(groups, keyword)
		}
	}
	slog.Info("matched groups", "groups", groups)
	return groups
}
