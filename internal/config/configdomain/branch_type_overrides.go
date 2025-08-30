package configdomain

import (
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
)

// BranchTypeOverrides contains all configured branch type overrides.
// These are stored in Git metadata like this: "git-town-branch.<name>.branchtype".
type BranchTypeOverrides map[gitdomain.LocalBranchName]BranchType

// adds the given BranchTypeOverrides to this BranchTypeOverrides
func (self BranchTypeOverrides) Concat(other BranchTypeOverrides) BranchTypeOverrides {
	result := make(BranchTypeOverrides, len(self)+len(other))
	for key, value := range self {
		result[key] = value
	}
	for key, value := range other {
		result[key] = value
	}
	return result
}
