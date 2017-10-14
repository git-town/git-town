package command

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/Originate/git-town/src/exit"
	"github.com/Originate/git-town/src/util"
)

// Command runs commands on the command line
type Command struct {
	Name   string
	Args   []string
	ran    bool
	err    error
	output string
}

// New creates a new Command instance
func New(command ...string) *Command {
	return &Command{Name: command[0], Args: command[1:]}
}

// Run runs this command.
// Doesn't run again if it ran already.
// Stores the outcome in fields of the instance.
func (r *Command) Run() {
	if r.ran {
		return
	}

	subProcess := exec.Command(r.Name, r.Args...) // #nosec
	output, err := subProcess.CombinedOutput()
	r.output = strings.TrimSpace(string(output))
	r.err = err
	r.ran = true
}

// RunOrExit runs this command.
// Doesn't run again if it ran already.
// Exits the application in case of errors
func (r *Command) RunOrExit() {
	r.Run()
	exit.OnWrapf(r.err, "Command: %s\nOutput: %s", r.String(), r.output)
}

// Output returns the output of this command.
// Runs if it hasn't so far.
func (r *Command) Output() string {
	r.Run()
	return r.output
}

// OutputOrExit returns the output of this command.
// Exits the application in case of errors
func (r *Command) OutputOrExit() string {
	r.RunOrExit()
	return r.output
}

// Err returns the error that this command encountered.
// Runs the command if it hasn't so far.
func (r *Command) Err() error {
	r.Run()
	return r.err
}

// OutputContainsLine returns whether the output of this command
// contains the given line
func (r *Command) OutputContainsLine(line string) bool {
	r.Run()
	lines := strings.Split(r.output, "\n")
	return util.DoesStringArrayContain(lines, line)
}

// OutputContainsText returns whether the output of this command
// contains the given text
func (r *Command) OutputContainsText(text string) bool {
	r.Run()
	return strings.Contains(r.output, text)
}

func (r Command) String() string {
	return fmt.Sprintf("%s %s", r.Name, strings.Join(r.Args, " "))
}
