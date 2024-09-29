package healthcheck

import (
	"fmt"
	"log/slog"
	"os/exec"
	"strconv"
	"strings"
)

// RepoSize checks the size of the repository to determine if it is oversized.
// It executes the 'git count-objects -v' command to obtain the size information,
func RepoSize() error {
	// TODO: reimplement here when go-git is available
	// https://github.com/go-git/go-git/blob/master/COMPATIBILITY.md
	cmd := exec.Command("git", "count-objects", "-v")
	output, err := cmd.CombinedOutput()
	if err != nil {
		slog.Error("running git command:", "err", err)
		return fmt.Errorf("error running git command: %w", err)
	}
	lines := strings.Split(string(output), "\n")
	var sum int
	for _, line := range lines {
		if strings.Contains(line, "size") {
			fields := strings.Fields(line)
			sizeStr := fields[1]
			size, err := strconv.Atoi(sizeStr)
			if err != nil {
				slog.Error("running git command:", "err", err)
				return fmt.Errorf("error running git command: %w", err)
			}
			sum += size
		}
	}
	if sum > 2048 {
		return fmt.Errorf("repository larger than 2MB. Please clean up or contact the teaching team.")
	}
	return nil
}
