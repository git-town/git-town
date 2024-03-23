package configdomain

import "github.com/git-town/git-town/v13/src/git/gitdomain"

type BranchType int

const (
	BranchTypeMainBranch BranchType = iota
	BranchTypePerennialBranch
	BranchTypeFeatureBranch
	BranchTypeParkedBranch
	BranchTypeContributionBranch
	BranchTypeObservedBranch
)

// ShouldPush indicates whether a branch with this type should push its local commit to origin.
func (self BranchType) ShouldPush(currentBranch, initialBranch gitdomain.LocalBranchName) bool {
	switch self {
	case BranchTypeMainBranch, BranchTypeFeatureBranch, BranchTypePerennialBranch, BranchTypeContributionBranch:
		return true
	case BranchTypeObservedBranch:
		return false
	case BranchTypeParkedBranch:
		return currentBranch == initialBranch
	}
	panic("unhandled branch type")
}

func (self BranchType) String() string {
	switch self {
	case BranchTypeMainBranch:
		return "main branch"
	case BranchTypePerennialBranch:
		return "perennial branch"
	case BranchTypeFeatureBranch:
		return "feature branch"
	case BranchTypeParkedBranch:
		return "parked branch"
	case BranchTypeContributionBranch:
		return "contribution branch"
	case BranchTypeObservedBranch:
		return "observed branch"
	}
	panic("unhandled branch type")
}
