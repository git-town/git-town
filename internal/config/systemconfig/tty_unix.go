//go:build !windows

package systemconfig

import "os"

func canOpenTTY() bool {
	f, err := os.Open("/dev/tty")
	if err != nil {
		return false
	}
	defer f.Close()
	return true
}
