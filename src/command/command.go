package command

import (
	"os/exec"
	"strings"

	"github.com/Originate/git-town/src/util"
	"github.com/acarl005/stripansi"
)

// Result contains the results of a command run in a subshell.
type Result struct {
	err    error
	output string
}

// Run executes the command given in argv notation.
func Run(argv ...string) *Result {
	name, args := argv[0], argv[1:]
	logRun(argv...)
	subProcess := exec.Command(name, args...) // #nosec
	output, err := subProcess.CombinedOutput()
	return &Result{err: err, output: stripansi.Strip(strings.TrimSpace(string(output)))}
}

// Output returns the output of this command.
// Runs if it hasn't so far.
func (c *Result) Output() string {
	return c.output
}

// OutputLines returns the output of this command, split into lines.
// Runs if it hasn't so far.
func (c *Result) OutputLines() []string {
	return strings.Split(c.Output(), "\n")
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
