package subshell

import (
	"strings"

	"github.com/acarl005/stripansi"
	"github.com/git-town/git-town/v7/src/stringslice"
)

// Result contains the results of a command run in a subshell.
type Output struct {
	Raw string // the raw output of the command
}

func NewOutput(raw []byte) *Output {
	return &Output{Raw: string(raw)}
}

// OutputLines provides the output of this command, split into lines.
// Runs if it hasn't so far.
func (c *Output) Lines() []string {
	return strings.Split(c.Sanitized(), "\n")
}

// OutputSanitized provides the output without ANSI color codes.
func (c *Output) Sanitized() string {
	return strings.TrimSpace(stripansi.Strip(c.Raw))
}

// OutputContainsLine returns whether the output of this command
// contains the given line.
func (c *Output) ContainsLine(line string) bool {
	return stringslice.Contains(c.Lines(), line)
}

// OutputContainsText returns whether the output of this command
// contains the given text.
func (c *Output) ContainsText(text string) bool {
	return strings.Contains(c.Raw, text)
}
