package healthcheck

import (
	"fmt"
	"os/exec"
)

func RepoLFS(rootDir string) error {
	cmd := exec.Command(
		"/usr/bin/git",
		"-c",
		"safe.directory=*",
		"lfs",
		"fsck",
		"--pointers",
	)
	cmd.Dir = rootDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf(
			"error running `git lfs fsck --pointers`: %w, output:\n%s",
			err,
			output,
		)
	}
	return nil
}
