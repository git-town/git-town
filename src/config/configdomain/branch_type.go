package configdomain

import (
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
)

type BranchType int

const (
	BranchTypeMainBranch BranchType = iota
	BranchTypePerennialBranch
	BranchTypeFeatureBranch
	BranchTypeParkedBranch
	BranchTypeContributionBranch
	BranchTypeObservedBranch
)

func NewBranchType(name string) Option[BranchType] {
	switch name {
	case "contribution":
		return Some(BranchTypeContributionBranch)
	case "feature":
		return Some(BranchTypeFeatureBranch)
	case "main":
		return Some(BranchTypeMainBranch)
	case "observed":
		return Some(BranchTypeObservedBranch)
	case "parked":
		return Some(BranchTypeParkedBranch)
	case "perennial":
		return Some(BranchTypePerennialBranch)
	case "(none)":
		return None[BranchType]()
	}
	panic("unknown branch type: " + name)
}

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
