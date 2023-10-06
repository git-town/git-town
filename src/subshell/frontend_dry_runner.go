package subshell

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/gohacks"
	"github.com/git-town/git-town/v9/src/messages"
)

// FrontendDryRunner prints the given shell commands to the CLI as if they were executed
// but does not execute them.
type FrontendDryRunner struct {
	GetCurrentBranch GetCurrentBranchFunc
	OmitBranchNames  bool
	CommandsCounter  *gohacks.Counter
}

// Run runs the given command in this ShellRunner's directory.
func (fdr *FrontendDryRunner) Run(executable string, args ...string) error {
	currentBranch, err := fdr.GetCurrentBranch()
	if err != nil {
		return err
	}
	PrintCommand(currentBranch, fdr.OmitBranchNames, executable, args...)
	return nil
}

// RunMany runs all given commands in current directory.
// Commands are provided as a list of argv-style strings.
// Failed commands abort immediately with the encountered error.
func (fdr *FrontendDryRunner) RunMany(commands [][]string) error {
	for _, argv := range commands {
		err := fdr.Run(argv[0], argv[1:]...)
		if err != nil {
			return fmt.Errorf(messages.RunCommandProblem, argv, err)
		}
	}
	return nil
}
