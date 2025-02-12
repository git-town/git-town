package configdomain

import (
	"fmt"
	"strings"

	. "github.com/git-town/git-town/v18/pkg/prelude"
)

type BranchType string

const (
	BranchTypeMainBranch         = BranchType("main")
	BranchTypePerennialBranch    = BranchType("perennial")
	BranchTypeFeatureBranch      = BranchType("feature")
	BranchTypeParkedBranch       = BranchType("parked")
	BranchTypeContributionBranch = BranchType("contribution")
	BranchTypeObservedBranch     = BranchType("observed")
	BranchTypePrototypeBranch    = BranchType("prototype")
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
	if len(text) == 0 || text == "(none)" {
		return None[BranchType](), nil
	}
	for _, branchType := range AllBranchTypes() {
		if strings.HasPrefix(branchType.String(), text) {
			return Some(branchType), nil
		}
	}
	return None[BranchType](), fmt.Errorf("unknown branch type: %q", text)
}

func (self BranchType) MustKnowParent() bool {
	switch self {
	case
		BranchTypeMainBranch,
		BranchTypePerennialBranch,
		BranchTypeContributionBranch,
		BranchTypeObservedBranch:
		return false
	case
		BranchTypeFeatureBranch,
		BranchTypeParkedBranch,
		BranchTypePrototypeBranch:
		return true
	}
	panic("unhandled branch type" + self.String())
}

// ShouldPush indicates whether a branch with this type should push its local commit to origin.
func (self BranchType) ShouldPush(isInitialBranch bool) bool {
	switch self {
	case
		BranchTypeMainBranch,
		BranchTypeFeatureBranch,
		BranchTypePerennialBranch,
		BranchTypeContributionBranch:
		return true
	case
		BranchTypeObservedBranch,
		BranchTypePrototypeBranch:
		return false
	case BranchTypeParkedBranch:
		return isInitialBranch
	}
	panic("unhandled branch type" + self.String())
}

func (self BranchType) String() string {
	return string(self)
}
