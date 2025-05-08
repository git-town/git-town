package opcodes

import (
	"github.com/git-town/git-town/v20/internal/config/configdomain"
	"github.com/git-town/git-town/v20/internal/git/gitdomain"
	"github.com/git-town/git-town/v20/internal/vm/shared"
	. "github.com/git-town/git-town/v20/pkg/prelude"
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
		opcodes := []shared.Opcode{}
		descendents := args.Config.Value.NormalConfig.Lineage.Descendants(self.BranchToDelete)
		if len(descendents) > 0 {
			opcodes = append(opcodes, &RebaseOntoRemoveDeleted{
				BranchToRebaseOnto: self.BranchToRebaseOnto,
				CommitsToRemove:    self.BranchToDelete.BranchName(),
				Upstream:           None[gitdomain.LocalBranchName](),
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
