package gitdomain

import (
	"fmt"
	"strings"
)

// BranchName is the name of a local or remote Git branch.
type BranchName string

func NewBranchName(id string) BranchName {
	if !isValidBranchName(id) {
		panic(fmt.Sprintf("%q is not a valid Git branch name", id))
	}
	return BranchName(id)
}

// IsLocal indicates whether the branch with this BranchName exists locally.
func (self BranchName) IsLocal() bool {
	return !strings.HasPrefix(string(self), "origin/")
}

// LocalName provides the local version of this branch name.
func (self BranchName) LocalName() LocalBranchName {
	return NewLocalBranchName(strings.TrimPrefix(string(self), "origin/"))
}

// RemoteName provides the remote version of this branch name.
func (self BranchName) RemoteName() RemoteBranchName {
	if strings.HasPrefix(string(self), "origin/") {
		return NewRemoteBranchName(string(self))
	}
	return NewRemoteBranchName("origin/" + string(self))
}

// Implementation of the fmt.Stringer interface.
func (self BranchName) String() string {
	return string(self)
}

// isValidBranchName indicates whether the given text is a valid branch name.
func isValidBranchName(text string) bool {
	return len(text) != 0
}
