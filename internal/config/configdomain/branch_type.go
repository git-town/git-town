package configdomain

import (
	"fmt"

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

func ParseBranchType(text string) (Option[BranchType], error) {
	switch text {
	case "contribution":
		return Some(BranchTypeContributionBranch), nil
	case "feature":
		return Some(BranchTypeFeatureBranch), nil
	case "main":
		return Some(BranchTypeMainBranch), nil
	case "observed":
		return Some(BranchTypeObservedBranch), nil
	case "parked":
		return Some(BranchTypeParkedBranch), nil
	case "perennial":
		return Some(BranchTypePerennialBranch), nil
	case "prototype":
		return Some(BranchTypePrototypeBranch), nil
	case "(none)":
		return None[BranchType](), nil
	}
	return None[BranchType](), fmt.Errorf("unknown branch type: %q", text)
}

func (self BranchType) MustKnowParent() bool {
	switch self {
	case BranchTypeMainBranch, BranchTypePerennialBranch, BranchTypeContributionBranch, BranchTypeObservedBranch:
		return false
	case BranchTypeFeatureBranch, BranchTypeParkedBranch, BranchTypePrototypeBranch:
		return true
	}
	panic("unhandled branch type" + self.String())
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
	panic("unhandled branch type" + self.String())
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
