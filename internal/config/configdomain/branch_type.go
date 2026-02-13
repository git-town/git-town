package configdomain

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v22/internal/messages"
	. "github.com/git-town/git-town/v22/pkg/prelude"
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

func ParseBranchType(text string, source string) (Option[BranchType], error) {
	if len(text) == 0 || text == messages.DialogResultNone {
		return None[BranchType](), nil
	}
	for _, branchType := range AllBranchTypes() {
		if strings.HasPrefix(branchType.String(), text) {
			return Some(branchType), nil
		}
	}
	return None[BranchType](), fmt.Errorf(messages.DialogResultUnknownBranchType, source, text)
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

// ShouldUpdateProposals indicates whether proposals for branches of this type
// should have their body updated with breadcrumbs.
// Observed and contribution branches cannot be proposed or shipped,
// so their proposals should not be modified by Git Town.
func (self BranchType) ShouldUpdateProposals() bool {
	switch self {
	case
		BranchTypeMainBranch,
		BranchTypePerennialBranch,
		BranchTypeFeatureBranch,
		BranchTypeParkedBranch,
		BranchTypePrototypeBranch:
		return true
	case
		BranchTypeObservedBranch,
		BranchTypeContributionBranch:
		return false
	}
	panic("unhandled branch type" + self.String())
}

func (self BranchType) String() string {
	return string(self)
}
