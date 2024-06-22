package subshell

import (
	"fmt"

	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/gohacks"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
)

// FrontendDryRunner prints the given shell commands to the CLI as if they were executed
// but does not execute them.
type FrontendDryRunner struct {
	Backend          gitdomain.Querier
	CommandsCounter  Mutable[gohacks.Counter]
	GetCurrentBranch GetCurrentBranchFunc
	OmitBranchNames  bool
	PrintCommands    bool
}

// Run runs the given command in this ShellRunner's directory.
func (self *FrontendDryRunner) Run(executable string, args ...string) error {
	var currentBranch gitdomain.LocalBranchName
	if self.OmitBranchNames {
		currentBranch = ""
	} else {
		var err error
		currentBranch, err = self.GetCurrentBranch(self.Backend)
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
