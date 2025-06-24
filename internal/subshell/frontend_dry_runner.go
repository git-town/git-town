package subshell

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/gohacks"
	"github.com/git-town/git-town/v21/internal/subshell/subshelldomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// FrontendDryRunner prints the given shell commands to the CLI as if they were executed
// but does not execute them.
type FrontendDryRunner struct {
	Backend          subshelldomain.Querier
	CommandsCounter  Mutable[gohacks.Counter]
	GetCurrentBranch GetCurrentBranchFunc
	PrintBranchNames bool
	PrintCommands    bool
}

func (self *FrontendDryRunner) Run(cmd string, args ...string) error {
	return self.execute([]string{}, cmd, args...)
}

func (self *FrontendDryRunner) RunWithEnv(env []string, cmd string, args ...string) error {
	return self.execute(env, cmd, args...)
}

// Run runs the given command in this ShellRunner's directory.
func (self *FrontendDryRunner) execute(env []string, executable string, args ...string) error {
	var currentBranch gitdomain.LocalBranchName
	if self.PrintBranchNames {
		var err error
		currentBranch, err = self.GetCurrentBranch(self.Backend)
		if err != nil {
			return err
		}
	}
	if self.PrintCommands {
		PrintCommand(currentBranch, self.PrintBranchNames, env, executable, args...)
		fmt.Println("(dry run)")
	}
	return nil
}
