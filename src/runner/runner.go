package runner

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/Originate/git-town/src/exit"
	"github.com/Originate/git-town/src/util"
)

// Runner runs commands on the command line
type Runner struct {
	Name   string
	Args   []string
	ran    bool
	err    error
	output string
}

// New creates a new Runner instance
func New(command ...string) *Runner {
	return &Runner{Name: command[0], Args: command[1:]}
}

// WithArg appends the given command-line switch to this runner's command
func (r *Runner) WithArg(arg string) *Runner {
	r.Args = append(r.Args, arg)
	return r
}

// WithArgs appends the given command-line switch to this runner's command
func (r *Runner) WithArgs(args ...string) *Runner {
	for _, arg := range args {
		r.WithArg(arg)
	}
	return r
}

// Run runs this runner.
// Doesn't run again if it ran already.
// Stores the outcome in fields of the instance.
func (r *Runner) Run() {
	if r.ran {
		return
	}

	subProcess := exec.Command(r.Name, r.Args...) // #nosec
	output, err := subProcess.CombinedOutput()
	r.output = strings.TrimSpace(string(output))
	r.err = err
	r.ran = true
}

// RunOrExit runs this runner.
// Doesn't run again if it ran already.
// Exits the application in case of errors
func (r *Runner) RunOrExit() {
	r.Run()
	exit.OnWrapf(r.err, "Command: %s\nOutput: %s", r.String(), r.output)
}

// Output returns the output of this command.
// Runs if it hasn't so far.
func (r *Runner) Output() string {
	r.Run()
	return r.output
}

// OutputOrExit returns the output of this command.
// Exits the application in case of errors
func (r *Runner) OutputOrExit() string {
	r.RunOrExit()
	return r.output
}

// Err returns the error that this runner encountered.
// Runs the command if it hasn't so far.
func (r *Runner) Err() error {
	r.Run()
	return r.err
}

// OutputContainsLine returns whether the output of this command
// contains the given line
func (r *Runner) OutputContainsLine(line string) bool {
	r.Run()
	lines := strings.Split(r.output, "\n")
	return util.DoesStringArrayContain(lines, line)
}

// OutputContainsText returns whether the output of this command
// contains the given text
func (r *Runner) OutputContainsText(text string) bool {
	r.Run()
	return strings.Contains(r.output, text)
}

func (r Runner) String() string {
	return fmt.Sprintf("%s %s", r.Name, strings.Join(r.Args, " "))
}
