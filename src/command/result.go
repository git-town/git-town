package command

import (
	"strings"

	"github.com/Originate/git-town/src/util"
	"github.com/acarl005/stripansi"
)

// Result contains the results of a command run in a subshell.
type Result struct {
	command string
	args    []string
	err     error
	output  string
}

// Args provides the arguments used when running the command.
func (c *Result) Args() []string {
	return c.args
}

// Command provides the command run that led to this result.
func (c *Result) Command() string {
	return c.command
}

// Output returns the output of this command.
// Runs if it hasn't so far.
func (c *Result) Output() string {
	return c.output
}

// OutputLines returns the output of this command, split into lines.
// Runs if it hasn't so far.
func (c *Result) OutputLines() []string {
	return strings.Split(c.OutputSanitized(), "\n")
}

// OutputSanitized provides the output without ANSI color codes.
func (c *Result) OutputSanitized() string {
	return strings.TrimSpace(stripansi.Strip(c.output))
}

// Err returns the error that this command encountered.
// Runs the command if it hasn't so far.
func (c *Result) Err() error {
	return c.err
}

// OutputContainsLine returns whether the output of this command
// contains the given line
func (c *Result) OutputContainsLine(line string) bool {
	return util.DoesStringArrayContain(c.OutputLines(), line)
}

// OutputContainsText returns whether the output of this command
// contains the given text
func (c *Result) OutputContainsText(text string) bool {
	return strings.Contains(c.output, text)
}
