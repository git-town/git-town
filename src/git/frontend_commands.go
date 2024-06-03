package git

import (
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
)

type FrontendRunner interface {
	Run(executable string, args ...string) error
	RunMany(commands [][]string) error
}

func (self *FrontendCommands) UndoLastCommit() error {
	return self.Runner.Run("git", "reset", "--soft", "HEAD~1")
}
