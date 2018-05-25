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
	ran    bool
	err    error
	output string
}

// New creates a new Command instance
func New(command ...string) *Command {
	return &Command{name: command[0], args: command[1:]}
}

// Run runs this command.
// Doesn't run again if it ran already.
// Stores the outcome in fields of the instance.
func (r *Command) Run() {
	if r.ran {
		return
	}

	logRun(r)
	subProcess := exec.Command(r.name, r.args...) // #nosec
	output, err := subProcess.CombinedOutput()
	r.output = stripansi.Strip(strings.TrimSpace(string(output)))
	r.err = err
	r.ran = true
}

// Output returns the output of this command.
// Runs if it hasn't so far.
func (r *Command) Output() string {
	r.Run()
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
	return fmt.Sprintf("%s %s", r.name, strings.Join(r.args, " "))
}
