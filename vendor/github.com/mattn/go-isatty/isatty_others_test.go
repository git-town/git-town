// +build !windows

package isatty

import (
	"os"
	"testing"
)

func TestCygwinPipeName(t *testing.T) {
	if IsCygwinTerminal(os.Stdout.Fd()) {
		t.Fatal("should be false always")
	}
}
