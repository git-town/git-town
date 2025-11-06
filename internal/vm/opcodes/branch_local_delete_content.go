package opcodes

import (
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/vm/shared"
)

// BranchLocalDeleteContent deletes the given branch including all commits.
type BranchLocalDeleteContent struct {
	BranchToDelete     gitdomain.LocalBranchName
	BranchToRebaseOnto gitdomain.LocalBranchName
}

func (self *BranchLocalDeleteContent) Run(args shared.RunArgs) error {
	switch args.Config.Value.NormalConfig.SyncFeatureStrategy {
	case configdomain.SyncFeatureStrategyRebase:
		opcodes := []shared.Opcode{}
		descendents := args.Config.Value.NormalConfig.Lineage.Descendants(self.BranchToDelete, args.Config.Value.NormalConfig.Order)
		if len(descendents) > 0 {
			opcodes = append(opcodes, &RebaseOntoRemoveDeleted{
				BranchToRebaseOnto: self.BranchToRebaseOnto,
				CommitsToRemove:    self.BranchToDelete.BranchName(),
			})
		}
		opcodes = append(opcodes, &BranchLocalDelete{
			Branch: self.BranchToDelete,
		})
		args.PrependOpcodes(opcodes...)
	case configdomain.SyncFeatureStrategyMerge, configdomain.SyncFeatureStrategyCompress:
		args.PrependOpcodes(
			&BranchLocalDelete{Branch: self.BranchToDelete},
		)
	}
	return nil
}
