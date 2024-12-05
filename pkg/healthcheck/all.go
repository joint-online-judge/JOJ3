package healthcheck

import (
	"fmt"

	"github.com/joint-online-judge/JOJ3/internal/conf"
)

type Result struct {
	Msg    string
	Failed bool
}

func All(
	confObj *conf.Conf,
	actor, repoName, rootDir, checkFileNameList, checkFileSumList string,
	groups, metaFile []string,
	repoSize float64,
) (res Result) {
	var err error
	if confObj != nil {
		output, err := TeapotCheck(confObj, actor, repoName, groups)
		if err != nil {
			res.Msg += fmt.Sprintf("### Teapot Check Failed:\n%s\n", output)
			res.Failed = true
		} else {
			res.Msg += fmt.Sprintf("### Teapot Check Result:\n%s\n", output)
		}
	}
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
	err = NonAsciiFiles(rootDir)
	if err != nil {
		res.Msg += fmt.Sprintf("### Non-ASCII Characters File Check Failed:\n%s\n", err.Error())
		res.Failed = true
	}
	err = NonAsciiMsg(rootDir)
	if err != nil {
		res.Msg += fmt.Sprintf("### Non-ASCII Characters Commit Message Check Failed:\n%s\n", err.Error())
		res.Failed = true
	}
	err = VerifyFiles(rootDir, checkFileNameList, checkFileSumList)
	if err != nil {
		res.Msg += fmt.Sprintf("### Repo File Check Failed:\n%s\n", err.Error())
		res.Failed = true
	}
	return
}
