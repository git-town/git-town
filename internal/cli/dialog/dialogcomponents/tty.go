package dialogcomponents

import (
	"errors"
	"os"

	"github.com/mattn/go-isatty"
)

// ErrNoTTY indicates that an interactive terminal is required but not available.
var ErrNoTTY = errors.New("no interactive terminal available")

// HasTTY reports whether an interactive terminal is available.
func HasTTY() bool {
	fd := os.Stdin.Fd()
	if isatty.IsTerminal(fd) || isatty.IsCygwinTerminal(fd) {
		return true
	}
	return canOpenTTY()
}

// RequireTTY returns ErrNoTTY when no interactive terminal is available.
func RequireTTY() error {
	if !HasTTY() {
		return ErrNoTTY
	}
	return nil
}
