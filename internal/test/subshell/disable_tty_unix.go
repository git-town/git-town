//go:build !windows

package subshell

import (
	"os/exec"
	"syscall"
)

func disableTTY(subProcess *exec.Cmd) {
	// HACK to work around a bug in the "exhaustruct" linter.
	// subProcess.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	sysProcAttr := new(syscall.SysProcAttr)
	sysProcAttr.Setsid = true
	subProcess.SysProcAttr = sysProcAttr
}
