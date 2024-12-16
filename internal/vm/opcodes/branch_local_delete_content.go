package opcodes

import (
	"github.com/git-town/git-town/v17/internal/config/configdomain"
	"github.com/git-town/git-town/v17/internal/git/gitdomain"
	"github.com/git-town/git-town/v17/internal/vm/shared"
	. "github.com/git-town/git-town/v17/pkg/prelude"
)

// deletes the given branch including all commits
type BranchLocalDeleteContent struct {
	BranchToDelete          gitdomain.LocalBranchName
	BranchToRebaseOnto      gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *BranchLocalDeleteContent) Run(args shared.RunArgs) error {
	switch args.Config.Value.NormalConfig.SyncFeatureStrategy {
	case configdomain.SyncFeatureStrategyRebase:
		args.PrependOpcodes(
			&RebaseOnto{
				BranchToRebaseAgainst: self.BranchToDelete.BranchName(),
				BranchToRebaseOnto:    self.BranchToRebaseOnto,
				Upstream:              None[gitdomain.LocalBranchName](),
			},
			&BranchLocalDelete{
				Branch: self.BranchToDelete,
			},
		)
	case configdomain.SyncFeatureStrategyMerge, configdomain.SyncFeatureStrategyCompress:
		args.PrependOpcodes(
			&BranchLocalDelete{Branch: self.BranchToDelete},
		)
	}
	return nil
}
