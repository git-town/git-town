package opcodes

import (
	"fmt"

	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/config/gitconfig"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/internal/vm/shared"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

var (
	branchTypeOverrideSet BranchTypeOverrideSet
	_                     shared.Runnable = &branchTypeOverrideSet
)

// BranchTypeOverrideSet registers the branch with the given name as a contribution branch in the Git config.
type BranchTypeOverrideSet struct {
	Branch     gitdomain.LocalBranchName
	BranchType configdomain.BranchType
}

func (self *BranchTypeOverrideSet) Run(args shared.RunArgs) error {
	if message, has := self.message().Get(); has {
		args.FinalMessages.Add(message)
	}
	return gitconfig.SetBranchTypeOverride(args.Backend, self.BranchType, self.Branch)
}

func (self *BranchTypeOverrideSet) message() Option[string] {
	switch self.BranchType {
	case configdomain.BranchTypeContributionBranch:
		return Some(fmt.Sprintf(messages.BranchIsNowContribution, self.Branch))
	case configdomain.BranchTypeFeatureBranch:
		return Some(fmt.Sprintf(messages.BranchIsNowFeature, self.Branch))
	case configdomain.BranchTypeMainBranch:
		return None[string]()
	case configdomain.BranchTypeObservedBranch:
		return Some(fmt.Sprintf(messages.BranchIsNowObserved, self.Branch))
	case configdomain.BranchTypeParkedBranch:
		return Some(fmt.Sprintf(messages.BranchIsNowParked, self.Branch))
	case configdomain.BranchTypePerennialBranch:
		return Some(fmt.Sprintf(messages.BranchIsNowPerennial, self.Branch))
	case configdomain.BranchTypePrototypeBranch:
		return Some(fmt.Sprintf(messages.BranchIsNowPrototype, self.Branch))
	}
	return None[string]()
}
