package opcodes

import (
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

// deletes the given branch including all commits
type BranchLocalDeleteContent struct {
	BranchToRebaseOnto      gitdomain.BranchName
	BranchToDelete          gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *BranchLocalDeleteContent) Run(args shared.RunArgs) error {
	switch args.Config.Value.NormalConfig.SyncFeatureStrategy {
	case configdomain.SyncFeatureStrategyRebase:
		args.PrependOpcodes(
			&RebaseOnto{BranchToRebaseOnto: self.BranchToRebaseOnto, BranchToRebaseAgainst: self.BranchToDelete},
			&BranchLocalDelete{Branch: self.BranchToDelete},
		)
	case configdomain.SyncFeatureStrategyMerge, configdomain.SyncFeatureStrategyCompress:
		args.PrependOpcodes(
			&BranchLocalDelete{Branch: self.BranchToDelete},
		)
	}
	return nil
}
