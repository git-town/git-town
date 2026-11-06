package systemconfig

import (
	"os"

	"github.com/git-town/git-town/v24/internal/config/configdomain"
	"github.com/mattn/go-isatty"
)

// DetermineStdioIsTerminal reports whether both STDIN and STDOUT are terminals,
// which subprocesses like gh need to prompt the user.
// Unlike DetermineTTY, it doesn't fall back to /dev/tty.
func DetermineStdioIsTerminal() configdomain.StdioIsTerminal {
	return configdomain.StdioIsTerminal(isTerminal(os.Stdin) && isTerminal(os.Stdout))
}

// DetermineTTY reports whether an interactive terminal is available.
func DetermineTTY() configdomain.HasTTY {
	if isTerminal(os.Stdin) {
		return true
	}
	return canOpenTTY()
}

func isTerminal(file *os.File) bool {
	fd := file.Fd()
	return isatty.IsTerminal(fd) || isatty.IsCygwinTerminal(fd)
}
