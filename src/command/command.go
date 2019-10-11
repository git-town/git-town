package command

import (
	"os"
	"os/exec"
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

// Run executes the command given in argv notation.
func Run(cmd string, args ...string) *Result {
	return RunInDir("", cmd, args...)
}

// RunInDir executes the given command in the given directory.
func RunInDir(dir string, cmd string, args ...string) *Result {
	return RunDirEnv(dir, os.Environ(), cmd, args...)
}

// RunDirEnv runs the given command in the given directory with the given environment.
func RunDirEnv(dir string, env []string, cmd string, args ...string) *Result {
	logRun(cmd, args...)
	subProcess := exec.Command(cmd, args...) // #nosec
	if dir != "" {
		subProcess.Dir = dir
	}
	subProcess.Env = env
	output, err := subProcess.CombinedOutput()
	return &Result{
		command: cmd,
		args:    args,
		err:     err,
		output:  stripansi.Strip(strings.TrimSpace(string(output))),
	}
}

// Args provids the arguments used when running the command.
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
