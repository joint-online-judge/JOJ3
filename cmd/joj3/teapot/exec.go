package teapot

import (
	"bufio"
	"bytes"
	"log/slog"
	"os/exec"
	"regexp"
	"sync"
)

func runCommand(args []string) (
	stdoutBuf *bytes.Buffer, err error,
) {
	re := regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
	cmd := exec.Command("joint-teapot", args...) // #nosec G204
	stdoutBuf = new(bytes.Buffer)
	cmd.Stdout = stdoutBuf
	stderr, err := cmd.StderrPipe()
	if err != nil {
		slog.Error("stderr pipe", "error", err)
		return
	}
	var wg sync.WaitGroup
	wg.Add(1)
	scanner := bufio.NewScanner(stderr)
	go func() {
		for scanner.Scan() {
			text := re.ReplaceAllString(scanner.Text(), "")
			if text == "" {
				continue
			}
			slog.Info("joint-teapot", "stderr", text)
		}
		wg.Done()
		if scanner.Err() != nil {
			slog.Error("stderr scanner", "error", scanner.Err())
		}
	}()
	if err = cmd.Start(); err != nil {
		slog.Error("cmd start", "error", err)
		return
	}
	wg.Wait()
	if err = cmd.Wait(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode := exitErr.ExitCode()
			slog.Error("cmd completed with non-zero exit code",
				"error", err,
				"exitCode", exitCode)
		} else {
			slog.Error("cmd wait", "error", err)
		}
		return
	}
	return
}
