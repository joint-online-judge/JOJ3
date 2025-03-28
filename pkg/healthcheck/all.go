// Package healthcheck provides a set of health checks for a repository.
package healthcheck

import (
	"fmt"
)

type Result struct {
	Msg    string
	Failed bool
}

func All(
	rootDir, checkFileNameList, checkFileSumList string,
	groups, metaFile []string,
	repoSize float64,
) (res Result) {
	var err error
	err = RepoSize(repoSize)
	if err != nil {
		res.Msg += fmt.Sprintf("### Repo Size Check Failed:\n%s\n", err.Error())
		res.Failed = true
	}
	err = ForbiddenCheck(rootDir)
	if err != nil {
		res.Msg += fmt.Sprintf("### Forbidden File Check Failed:\n%s\n", err.Error())
		res.Failed = true
	}
	err = MetaCheck(rootDir, metaFile)
	if err != nil {
		res.Msg += fmt.Sprintf("### Meta File Check Failed:\n%s\n", err.Error())
		res.Failed = true
	}
	err = NonASCIIFiles(rootDir)
	if err != nil {
		res.Msg += fmt.Sprintf("### Non-ASCII Characters File Check Failed:\n%s\n", err.Error())
		res.Failed = true
	}
	err = NonASCIIMsg(rootDir)
	if err != nil {
		res.Msg += fmt.Sprintf("### Non-ASCII Characters Commit Message Check Failed:\n%s\n", err.Error())
		res.Failed = true
	}
	err = VerifyFiles(rootDir, checkFileNameList, checkFileSumList)
	if err != nil {
		res.Msg += fmt.Sprintf("### Repo File Check Failed:\n%s\n", err.Error())
		res.Failed = true
	}
	return res
}
