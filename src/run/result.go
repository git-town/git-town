package run

import (
	"fmt"
	"strings"

	"github.com/acarl005/stripansi"
	"github.com/git-town/git-town/v7/src/stringslice"
)

// Result contains the results of a command run in a subshell.
type Result struct {
	args     []string // arguments for the executed command
	command  string   // the executed command
	exitCode int      // the exit code of the command
	output   string   // the raw output of the command
}

// Args provides the arguments used when running the command.
func (c *Result) Args() []string {
	return c.args
}

// Command provides the command run that led to this result.
func (c *Result) Command() string {
	return c.command
}

// ExitCode provides the exit code of the command.
func (c *Result) ExitCode() int {
	return c.exitCode
}

// FullCmd provides the full command run.
func (c *Result) FullCmd() string {
	return fmt.Sprintf("%s %s", c.command, strings.Join(c.args, " "))
}

// Output provides the output of this command.
// Runs if it hasn't so far.
func (c *Result) Output() string {
	return c.output
}

// OutputLines provides the output of this command, split into lines.
// Runs if it hasn't so far.
func (c *Result) OutputLines() []string {
	return strings.Split(c.OutputSanitized(), "\n")
}

// OutputSanitized provides the output without ANSI color codes.
func (c *Result) OutputSanitized() string {
	return strings.TrimSpace(stripansi.Strip(c.output))
}

// OutputContainsLine returns whether the output of this command
// contains the given line.
func (c *Result) OutputContainsLine(line string) bool {
	return stringslice.Contains(c.OutputLines(), line)
}

// OutputContainsText returns whether the output of this command
// contains the given text.
func (c *Result) OutputContainsText(text string) bool {
	return strings.Contains(c.output, text)
}
