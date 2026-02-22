//go:build windows

package subshell

import (
	"os/exec"
)

func disableTTY(subProcess *exec.Cmd) {
	panic("disabling TTY on Windows not implemented yet")
}
