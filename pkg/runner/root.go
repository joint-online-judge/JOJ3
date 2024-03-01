package runner

import (
	"bytes"
	"log/slog"
	"os/exec"
	"os/user"
	"strconv"
	"syscall"
	"time"

	"github.com/containerd/cgroups"
	"github.com/opencontainers/runtime-spec/specs-go"
)

type RunResult struct {
	ReturnCode int
	Stdout     []byte
	Stderr     []byte
	TimedOut   bool
	TimeNs     uint64
	MemoryByte uint64
}

func RunInCgroupsV1(
	args []string, username string, cgroupsPath string, timeoutMs uint,
) (result *RunResult, err error) {
	u, err := user.Lookup(username)
	if err != nil {
		return
	}
	control, err := cgroups.New(
		cgroups.V1,
		cgroups.StaticPath(cgroupsPath),
		&specs.LinuxResources{},
	)
	if err != nil {
		return
	}
	defer func() {
		if err := control.Delete(); err != nil {
			slog.Error("control.Delete", "error", err)
		}
	}()
	cmd := exec.Command(args[0], args[1:]...)
	cmd.SysProcAttr = &syscall.SysProcAttr{}
	uid, _ := strconv.ParseUint(u.Uid, 10, 32)
	gid, _ := strconv.ParseUint(u.Gid, 10, 32)
	cmd.SysProcAttr.Credential = &syscall.Credential{Uid: uint32(uid), Gid: uint32(gid)}
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	start := time.Now()
	err = cmd.Start()
	if err != nil {
		return
	}
	pid := cmd.Process.Pid
	slog.Info("process created", "pid", strconv.Itoa(pid))
	if err = control.Add(cgroups.Process{Pid: pid}); err != nil {
		return
	}
	var returnCode int
	exitCode := make(chan int, 1)
	go func(exit_code chan int) {
		if err = cmd.Wait(); err != nil {
			exit_code <- err.(*exec.ExitError).ExitCode()
		} else {
			exit_code <- 0
		}
	}(exitCode)
	timeoutLimit := time.Duration(timeoutMs) * time.Millisecond
	timedOut := false
	select {
	case returnCode = <-exitCode:
		slog.Info("done success", "time", time.Since(start))
	case <-time.After(timeoutLimit):
		slog.Info("done timeout", "time", time.Since(start))
		_ = cmd.Process.Kill()
		returnCode = <-exitCode
		timedOut = true
	}
	stats, err := control.Stat(cgroups.IgnoreNotExist)
	if err != nil {
		return
	}
	result = &RunResult{
		ReturnCode: returnCode,
		Stdout:     stdout.Bytes(),
		Stderr:     stderr.Bytes(),
		TimedOut:   timedOut,
		TimeNs:     stats.CPU.Usage.Kernel + stats.CPU.Usage.User,
		MemoryByte: stats.Memory.Usage.Max, // Memory.Usage.Max = 0 when killed
	}
	return
}
