//go:build !windows

package subshell

import (
	"os/exec"
	"syscall"
)

func disableTTY(subProcess *exec.Cmd) {
	subProcess.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
}
