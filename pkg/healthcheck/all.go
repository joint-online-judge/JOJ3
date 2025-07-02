// Package healthcheck provides a set of health checks for a repository.
package healthcheck

import (
	"fmt"
	"strings"
)

type Result struct {
	Msg    string
	Failed bool
}

func All(
	rootDir, checkFileNameList, checkFileSumList, allowedDomainList, actorCsvPath string,
	metaFile []string, repoSize float64,
) (res Result) {
	var err error
	err = RepoSize(rootDir, repoSize)
	if err != nil {
		res.Msg += fmt.Sprintf("### Repo Size Check Failed:\n%s\n", err.Error())
		res.Failed = true
	} else {
		res.Msg += "### Repo Size Check Passed\n"
	}
	err = RepoLFS(rootDir)
	if err != nil {
		res.Msg += fmt.Sprintf("### Repo LFS Check Failed:\n%s\n", err.Error())
		res.Failed = true
	} else {
		res.Msg += "### Repo LFS Check Passed\n"
	}
	err = ForbiddenCheck(rootDir)
	if err != nil {
		res.Msg += fmt.Sprintf("### Forbidden File Check Failed:\n%s\n", err.Error())
		res.Failed = true
	} else {
		res.Msg += "### Forbidden File Check Passed\n"
	}
	err = MetaCheck(rootDir, metaFile)
	if err != nil {
		res.Msg += fmt.Sprintf("### Meta File Check Failed:\n%s\n", err.Error())
		res.Failed = true
	} else {
		res.Msg += "### Meta File Check Passed\n"
	}
	err = NonASCIIFiles(rootDir)
	if err != nil {
		res.Msg += fmt.Sprintf("### Non-ASCII Characters File Check Failed:\n%s\n", err.Error())
		res.Failed = true
	} else {
		res.Msg += "### Non-ASCII Characters File Check Passed\n"
	}
	err = NonASCIIMsg(rootDir)
	if err != nil {
		res.Msg += fmt.Sprintf("### Non-ASCII Characters Commit Message Check Failed:\n%s\n", err.Error())
		res.Failed = true
	} else {
		res.Msg += "### Non-ASCII Characters Commit Message Check Passed\n"
	}
	err = VerifyFiles(rootDir, checkFileNameList, checkFileSumList)
	if err != nil {
		res.Msg += fmt.Sprintf("### Repo File Check Failed:\n%s\n", err.Error())
		res.Failed = true
	} else {
		res.Msg += "### Repo File Check Passed\n"
	}
	err = AuthorEmailCheck(rootDir, strings.Split(allowedDomainList, ","), actorCsvPath)
	if err != nil {
		res.Msg += fmt.Sprintf("### Author Email Check Failed:\n%s\n", err.Error())
		res.Failed = true
	} else {
		res.Msg += "### Author Email Check Passed\n"
	}
	return res
}
