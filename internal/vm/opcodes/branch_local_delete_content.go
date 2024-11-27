package opcodes

import (
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

// deletes the given branch including all commits
type BranchLocalDeleteContent struct {
	BranchToDelete          gitdomain.LocalBranchName
	BranchToRebaseOnto      gitdomain.BranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *BranchLocalDeleteContent) Run(args shared.RunArgs) error {
	switch args.Config.Value.NormalConfig.SyncFeatureStrategy {
	case configdomain.SyncFeatureStrategyRebase:
		args.PrependOpcodes(
			&RebaseOnto{
				BranchToRebaseAgainst: self.BranchToDelete,
				BranchToRebaseOnto:    self.BranchToRebaseOnto,
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
