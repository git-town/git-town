package systemconfig

import (
	"os"

	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/mattn/go-isatty"
)

// HasTTY reports whether an interactive terminal is available.
func DetermineTTY() configdomain.HasTTY {
	fd := os.Stdin.Fd()
	if isatty.IsTerminal(fd) || isatty.IsCygwinTerminal(fd) {
		return true
	}
	return canOpenTTY()
}
