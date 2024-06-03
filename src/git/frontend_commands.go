package git

import (
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/config/gitconfig"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
)

type FrontendRunner interface {
	Run(executable string, args ...string) error
	RunMany(commands [][]string) error
}

// SetHostingPlatform sets the given code hosting platform.
func (self *FrontendCommands) SetOriginHostname(hostname configdomain.HostingOriginHostname) error {
	return self.Runner.Run("git", "config", gitconfig.KeyHostingOriginHostname.String(), hostname.String())
}

// SquashMerge squash-merges the given branch into the current branch.
func (self *FrontendCommands) SquashMerge(branch gitdomain.LocalBranchName) error {
	return self.Runner.Run("git", "merge", "--squash", "--ff", branch.String())
}

// StageFiles adds the file with the given name to the Git index.
func (self *FrontendCommands) StageFiles(names ...string) error {
	args := append([]string{"add"}, names...)
	return self.Runner.Run("git", args...)
}

// StartCommit starts a commit and stops at asking the user for the commit message.
func (self *FrontendCommands) StartCommit() error {
	return self.Runner.Run("git", "commit")
}

// Stash adds the current files to the Git stash.
func (self *FrontendCommands) Stash() error {
	return self.Runner.RunMany([][]string{
		{"git", "add", "-A"},
		{"git", "stash"},
	})
}

func (self *FrontendCommands) UndoLastCommit() error {
	return self.Runner.Run("git", "reset", "--soft", "HEAD~1")
}
