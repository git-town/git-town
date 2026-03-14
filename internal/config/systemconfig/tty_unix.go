//go:build !windows

package systemconfig

import (
	"os"

	"github.com/git-town/git-town/v22/internal/config/configdomain"
)

func canOpenTTY() configdomain.HasTTY {
	f, err := os.Open("/dev/tty")
	if err != nil {
		return false
	}
	defer f.Close()
	return true
}
