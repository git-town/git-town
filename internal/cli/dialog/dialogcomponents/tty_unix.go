//go:build !windows

package dialogcomponents

import "os"

func canOpenTTY() bool {
	f, err := os.Open("/dev/tty")
	if err != nil {
		return false
	}
	defer f.Close()
	return true
}
