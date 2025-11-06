package gitdomain

import (
	"fmt"
	"strings"

	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// RemoteBranchName is the name of a remote branch, e.g. "origin/foo".
type RemoteBranchName string

func NewRemoteBranchName(id string) RemoteBranchName {
	if !isValidRemoteBranchName(id) {
		panic(fmt.Sprintf("%q is not a valid remote branch name", id))
	}
	return RemoteBranchName(id)
}

func NewRemoteBranchNameOption(id string) Option[RemoteBranchName] {
	if isValidRemoteBranchName(id) {
		return Some(NewRemoteBranchName(id))
	}
	return None[RemoteBranchName]()
}

func isValidRemoteBranchName(value string) bool {
	if len(value) < 3 {
		return false
	}
	if !strings.Contains(value, "/") {
		return false
	}
	return true
}

// BranchName widens the type of this RemoteBranchName to a more generic BranchName.
func (self RemoteBranchName) BranchName() BranchName {
	return BranchName(string(self))
}

// LocalBranchName provides the name of the local branch that this remote branch tracks.
func (self RemoteBranchName) LocalBranchName() LocalBranchName {
	_, localBranch := self.Parts()
	return localBranch
}

func (self RemoteBranchName) Parts() (Remote, LocalBranchName) {
	remoteName, branchname, _ := strings.Cut(string(self), "/")
	return Remote(remoteName), NewLocalBranchName(branchname)
}

func (self RemoteBranchName) Remote() Remote {
	remote, _ := self.Parts()
	return remote
}

func (self RemoteBranchName) String() string {
	return string(self)
}
