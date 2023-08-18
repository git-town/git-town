package domain

import (
	"fmt"
	"strings"
)

// BranchName is the name of a local or remote Git branch.
type BranchName struct {
	Location // a BranchName is a type of Location
}

func NewBranchName(value string) BranchName {
	if !isValidBranchName(value) {
		panic(fmt.Sprintf("%q is not a valid Git branch name", value))
	}
	return BranchName{Location{value}}
}

// IsLocal indicates whether the branch with this BranchName exists locally.
func (c BranchName) IsLocal() bool {
	return !strings.HasPrefix(c.id, "origin/")
}

// LocalName provides the local version of this branch name.
func (c BranchName) LocalName() LocalBranchName {
	return NewLocalBranchName(strings.TrimPrefix(c.id, "origin/"))
}

// RemoteName provides the remote version of this branch name.
func (c BranchName) RemoteName() RemoteBranchName {
	if strings.HasPrefix(c.id, "origin/") {
		return NewRemoteBranchName(c.id)
	}
	return NewRemoteBranchName("origin/" + c.id)
}

// Implementation of the fmt.Stringer interface.
func (c BranchName) String() string { return c.id }

// isValidBranchName indicates whether the given text is a valid branch name.
func isValidBranchName(text string) bool {
	return len(text) != 0
}
