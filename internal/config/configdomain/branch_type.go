package configdomain

import (
	. "github.com/git-town/git-town/v16/pkg/prelude"
)

type BranchType int

const (
	BranchTypeMainBranch BranchType = iota
	BranchTypePerennialBranch
	BranchTypeFeatureBranch
	BranchTypeParkedBranch
	BranchTypeContributionBranch
	BranchTypeObservedBranch
	BranchTypePrototypeBranch
)

func ParseBranchType(name string) Option[BranchType] {
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
	case "prototype":
		return Some(BranchTypePrototypeBranch)
	case "(none)":
		return None[BranchType]()
	}
	panic("unknown branch type: " + name)
}

// ShouldPush indicates whether a branch with this type should push its local commit to origin.
func (self BranchType) ShouldPush(isInitialBranch bool) bool {
	switch self {
	case BranchTypeMainBranch, BranchTypeFeatureBranch, BranchTypePerennialBranch, BranchTypeContributionBranch:
		return true
	case BranchTypeObservedBranch, BranchTypePrototypeBranch:
		return false
	case BranchTypeParkedBranch:
		return isInitialBranch
	}
	panic("unhandled branch type")
}

func (self BranchType) String() string {
	switch self {
	case BranchTypeContributionBranch:
		return "contribution branch"
	case BranchTypeFeatureBranch:
		return "feature branch"
	case BranchTypeMainBranch:
		return "main branch"
	case BranchTypeObservedBranch:
		return "observed branch"
	case BranchTypeParkedBranch:
		return "parked branch"
	case BranchTypePerennialBranch:
		return "perennial branch"
	case BranchTypePrototypeBranch:
		return "prototype branch"
	}
	panic("unhandled branch type")
}
