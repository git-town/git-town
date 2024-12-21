package undodomain

import "github.com/git-town/git-town/v17/internal/git/gitdomain"

// InconsistentChange describes a change where both local and remote branch exist before and after,
// but it's not an OmniChange, i.e. the SHA are different.
type InconsistentChange struct {
	Before gitdomain.BranchInfo
	After  gitdomain.BranchInfo
}
