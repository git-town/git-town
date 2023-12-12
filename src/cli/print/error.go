package print

import (
	"github.com/fatih/color"
	"github.com/git-town/git-town/v11/src/cli/io"
)

// Error prints the given error message to the console.
func Error(err error) {
	io.PrintlnColor(color.New(color.Bold).Add(color.FgRed), "\nError:", err.Error(), "\n")
}
