package subshell

import (
	"fmt"

	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/gohacks"
	"github.com/git-town/git-town/v14/src/messages"
)

// FrontendDryRunner prints the given shell commands to the CLI as if they were executed
// but does not execute them.
type FrontendDryRunner struct {
	CommandsCounter  gohacks.Counter
	GetCurrentBranch GetCurrentBranchFunc
	OmitBranchNames  bool
	PrintCommands    bool
}

// Run runs the given command in this ShellRunner's directory.
func (self *FrontendDryRunner) Run(executable string, args ...string) error {
	var currentBranch gitdomain.LocalBranchName
	if self.OmitBranchNames {
		currentBranch = gitdomain.EmptyLocalBranchName()
	} else {
		var err error
		currentBranch, err = self.GetCurrentBranch()
		if err != nil {
			return err
		}
	}
	if self.PrintCommands {
		PrintCommand(currentBranch, self.OmitBranchNames, executable, args...)
		fmt.Println("(dry run)")
	}
	return nil
}

// RunMany runs all given commands in current directory.
// Commands are provided as a list of argv-style strings.
// Failed commands abort immediately with the encountered error.
func (self *FrontendDryRunner) RunMany(commands [][]string) error {
	for _, argv := range commands {
		err := self.Run(argv[0], argv[1:]...)
		if err != nil {
			return fmt.Errorf(messages.RunCommandProblem, argv, err)
		}
	}
	return nil
}
