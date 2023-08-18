package domain

import (
	"fmt"
	"strings"
)

// BranchName is a branch name as Git Town encounters it. It could be either a LocalBranchName or a RemoteBranchName.
// Local branches have the BranchName "foo". Remote branches have the BranchName "origin/foo".
type BranchName struct {
	Location
}

// Implements the fmt.Stringer interface.
func (c BranchName) String() string { return c.id }

// IsLocal indicates whether the branch with this BranchName exists locally.
func (c BranchName) IsLocal() bool {
	return !strings.HasPrefix(c.id, "origin/")
}

// LocalName provides the name that a branch with the given BranchName would have if it was local.
func (c BranchName) LocalName() LocalBranchName {
	return NewLocalBranchName(strings.TrimPrefix(c.id, "origin/"))
}

// RemoteName provides the name that a branch with the given BranchName would have if it was remote.
func (c BranchName) RemoteName() RemoteBranchName {
	if strings.HasPrefix(c.id, "origin/") {
		return NewRemoteBranchName(c.id)
	}
	return NewRemoteBranchName("origin/" + c.id)
}

func NewBranchName(value string) BranchName {
	if !isValidBranchName(value) {
		panic(fmt.Sprintf("%q is not a valid Git branch name", value))
	}
	return BranchName{Location{value}}
}

func isValidBranchName(value string) bool {
	return len(value) != 0
}

type BranchNames []BranchName

func (b BranchNames) Strings() []string {
	result := make([]string, len(b))
	for b, branch := range b {
		result[b] = branch.id
	}
	return result
}
