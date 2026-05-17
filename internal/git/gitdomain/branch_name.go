package gitdomain

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v23/internal/gohacks/stringss"
)

// BranchName is the name of a local or remote Git branch.
type BranchName stringss.Trimmed

func BranchNameOrPanic(id stringss.Trimmed) BranchName {
	if !isValidBranchName(id) {
		panic(fmt.Sprintf("%q is not a valid Git branch name", id))
	}
	return BranchName(id)
}

// IsLocal indicates whether the branch with this BranchName exists locally.
func (self BranchName) IsLocal() bool {
	return !strings.HasPrefix(string(self), "origin/")
}

// LocalName provides the (theoretical) local version of this branch name.
func (self BranchName) LocalName() LocalBranchName {
	return LocalBranchNameOrPanic(stringss.Trim(strings.TrimPrefix(string(self), "origin/")))
}

func (self BranchName) Location() Location {
	return NewLocation(self.String())
}

// RefName provides the fully qualified reference name for this branch.
func (self BranchName) RefName() string {
	if self.IsLocal() {
		return "refs/heads/" + self.String()
	}
	return self.String()
}

// RemoteName provides the remote version of this branch name.
func (self BranchName) RemoteName() RemoteBranchName {
	if strings.HasPrefix(string(self), "origin/") {
		return RemoteBranchNameOrPanic(string(self))
	}
	return RemoteBranchNameOrPanic("origin/" + string(self))
}

// String implements the fmt.Stringer interface.
func (self BranchName) String() string {
	return string(self)
}

// isValidBranchName indicates whether the given text is a valid branch name.
func isValidBranchName(text stringss.Trimmed) bool {
	return len(text) != 0
}
