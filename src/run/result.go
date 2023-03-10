package run

import (
	"fmt"
	"strings"

	"github.com/acarl005/stripansi"
	"github.com/git-town/git-town/v7/src/stringslice"
)

// Result contains the results of a command run in a subshell.
type Result struct {
	Args     []string // arguments for the executed command
	Command  string   // the executed command
	ExitCode int      // the exit code of the command
	Output   string   // the raw output of the command
}

// FullCmd provides the full command run.
func (c *Result) FullCmd() string {
	return fmt.Sprintf("%s %s", c.Command, strings.Join(c.Args, " "))
}

// OutputLines provides the output of this command, split into lines.
// Runs if it hasn't so far.
func (c *Result) OutputLines() []string {
	return strings.Split(c.OutputSanitized(), "\n")
}

// OutputSanitized provides the output without ANSI color codes.
func (c *Result) OutputSanitized() string {
	return strings.TrimSpace(stripansi.Strip(c.Output))
}

// OutputContainsLine returns whether the output of this command
// contains the given line.
func (c *Result) OutputContainsLine(line string) bool {
	return stringslice.Contains(c.OutputLines(), line)
}

// OutputContainsText returns whether the output of this command
// contains the given text.
func (c *Result) OutputContainsText(text string) bool {
	return strings.Contains(c.Output, text)
}
