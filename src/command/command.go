package command

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/Originate/git-town/src/util"
	"github.com/acarl005/stripansi"
)

// Command runs commands on the command line
type Command struct {
	name   string
	args   []string
	dir    string
	ran    bool
	err    error
	output string
}

// New creates a new Command instance
func New(command ...string) *Command {
	return &Command{name: command[0], args: command[1:]}
}

// NewInDir creates a new Command instance that runs in the given directory.
func NewInDir(dir string, command ...string) *Command {
	result := New(command...)
	result.dir = dir
	return result
}

// Dir provides the directory in which this command is supposed to run.
func (c *Command) Dir() string {
	return c.dir
}

// Run runs this command.
// Doesn't run again if it ran already.
// Stores the outcome in fields of the instance.
func (c *Command) Run() {
	if c.ran {
		return
	}

	logRun(c)
	subProcess := exec.Command(c.name, c.args...) // #nosec
	if c.dir != "" {
		subProcess.Dir = c.dir
	}
	output, err := subProcess.CombinedOutput()
	c.output = stripansi.Strip(strings.TrimSpace(string(output)))
	c.err = err
	c.ran = true
}

// Output returns the output of this command.
// Runs if it hasn't so far.
func (c *Command) Output() string {
	c.Run()
	return c.output
}

// OutputLines returns the output of this command, split into lines.
// Runs if it hasn't so far.
func (c *Command) OutputLines() []string {
	return strings.Split(c.Output(), "\n")
}

// Err returns the error that this command encountered.
// Runs the command if it hasn't so far.
func (c *Command) Err() error {
	c.Run()
	return c.err
}

// OutputContainsLine returns whether the output of this command
// contains the given line
func (c *Command) OutputContainsLine(line string) bool {
	return util.DoesStringArrayContain(c.OutputLines(), line)
}

// OutputContainsText returns whether the output of this command
// contains the given text
func (c *Command) OutputContainsText(text string) bool {
	c.Run()
	return strings.Contains(c.output, text)
}

func (c Command) String() string {
	return fmt.Sprintf("%s %s", c.name, strings.Join(c.args, " "))
}
