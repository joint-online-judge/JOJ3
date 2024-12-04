package local

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/criyle/go-judge/envexec"
	"github.com/joint-online-judge/JOJ3/internal/stage"
)

type Local struct{}

func (e *Local) Run(cmds []stage.Cmd) ([]stage.ExecutorResult, error) {
	var results []stage.ExecutorResult

	for _, cmd := range cmds {
		execCmd := exec.Command(cmd.Args[0], cmd.Args[1:]...) // #nosec G204

		if len(cmd.Env) > 0 {
			env := os.Environ()
			env = append(env, cmd.Env...)
			execCmd.Env = env
		}

		if cmd.Stdin != nil {
			if cmd.Stdin.Content != nil {
				execCmd.Stdin = strings.NewReader(*cmd.Stdin.Content)
			} else if cmd.Stdin.Src != nil {
				file, err := os.Open(*cmd.Stdin.Src)
				if err != nil {
					return nil, fmt.Errorf("failed to open stdin file: %v", err)
				}
				defer file.Close()
				execCmd.Stdin = file
			}
		}
		var stdoutBuffer, stderrBuffer bytes.Buffer
		execCmd.Stdout = &stdoutBuffer
		execCmd.Stderr = &stderrBuffer

		startTime := time.Now()
		err := execCmd.Start()
		if err != nil {
			return nil, fmt.Errorf("failed to start command: %v", err)
		}

		err = execCmd.Wait()
		endTime := time.Now()
		runTime := endTime.Sub(startTime)
		processState := execCmd.ProcessState
		result := stage.ExecutorResult{
			Status:     stage.Status(envexec.StatusAccepted),
			ExitStatus: processState.ExitCode(),
			Error:      "",
			Time: func() uint64 {
				nanos := processState.UserTime().Nanoseconds()
				if nanos < 0 {
					return 0
				}
				return uint64(nanos)
			}(),
			Memory: func() uint64 {
				usage := processState.SysUsage()
				rusage, ok := usage.(*syscall.Rusage)
				if !ok {
					return 0
				}
				maxRssKB := rusage.Maxrss
				maxRssBytes := maxRssKB * 1024
				if maxRssBytes < 0 {
					return 0
				}
				return uint64(maxRssBytes)
			}(),
			RunTime: func() uint64 {
				nanos := runTime.Nanoseconds()
				if nanos < 0 {
					return 0
				}
				return uint64(nanos)
			}(),
			Files:   map[string]string{},
			FileIDs: map[string]string{},
		}

		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				result.Status = stage.Status(envexec.StatusNonzeroExitStatus)
				result.Error = exitErr.Error()
			} else {
				result.Status = stage.Status(envexec.StatusInternalError)
				result.Error = err.Error()
			}
		}

		if cmd.Stdout != nil && cmd.Stdout.Name != nil {
			result.Files[*cmd.Stdout.Name] = stdoutBuffer.String()
		}
		if cmd.Stderr != nil && cmd.Stderr.Name != nil {
			result.Files[*cmd.Stderr.Name] = stderrBuffer.String()
		}

		if err := handleCopyOut(&result, cmd); err != nil {
			result.Status = stage.Status(envexec.StatusFileError)
			result.Error = err.Error()
		}

		results = append(results, result)
	}

	return results, nil
}

// Helper function to handle copyOut files
func handleCopyOut(result *stage.ExecutorResult, cmd stage.Cmd) error {
	for _, filename := range cmd.CopyOut {
		if _, ok := result.Files[filename]; ok {
			continue
		}
		optional := false
		if strings.HasSuffix(filename, "?") {
			optional = true
			filename = strings.TrimSuffix(filename, "?")
		}
		result.Files[filename] = ""
		// Read file and add to result.Files
		file, err := os.Open(filename)
		if err != nil {
			if !optional {
				return err
			}
			continue
		}
		defer file.Close()
		content, err := io.ReadAll(file)
		if err != nil {
			return err
		}
		result.Files[filename] = string(content)
	}
	return nil
}

func (e *Local) Cleanup() error {
	return nil
}
