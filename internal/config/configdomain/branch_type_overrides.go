package configdomain

import (
	"maps"

	"github.com/git-town/git-town/v22/internal/git/gitdomain"
)

// BranchTypeOverrides contains all configured branch type overrides.
// These are stored in Git metadata like this: "git-town-branch.<name>.branchtype".
type BranchTypeOverrides map[gitdomain.LocalBranchName]BranchType

// Concat adds the given BranchTypeOverrides to this BranchTypeOverrides.
func (self BranchTypeOverrides) Concat(other BranchTypeOverrides) BranchTypeOverrides {
	result := make(BranchTypeOverrides, len(self)+len(other))
	maps.Copy(result, self)
	maps.Copy(result, other)
	return result
}
