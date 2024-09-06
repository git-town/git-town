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

func AllBranchTypes() []BranchType {
	return []BranchType{
		BranchTypeMainBranch,
		BranchTypePerennialBranch,
		BranchTypeFeatureBranch,
		BranchTypeParkedBranch,
		BranchTypeContributionBranch,
		BranchTypeObservedBranch,
		BranchTypePrototypeBranch,
	}
}

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
	case "(none)", "":
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
		return "contribution"
	case BranchTypeFeatureBranch:
		return "feature"
	case BranchTypeMainBranch:
		return "main"
	case BranchTypeObservedBranch:
		return "observed"
	case BranchTypeParkedBranch:
		return "parked"
	case BranchTypePerennialBranch:
		return "perennial"
	case BranchTypePrototypeBranch:
		return "prototype"
	}
	panic("unhandled branch type")
}
