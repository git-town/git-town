package sync

import (
	"github.com/git-town/git-town/v22/internal/config"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/vm/opcodes"
	"github.com/git-town/git-town/v22/internal/vm/program"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/git-town/git-town/v22/pkg/set"
)

type AddOpcodesToUpdateProposalStackArgs struct {
	ChangedBranches gitdomain.LocalBranchNames // all branches that were modified in the current command
	Config          config.ValidatedConfig
	Program         Mutable[program.Program]
}

func AddOpcodesToUpdateProposalStack(args AddOpcodesToUpdateProposalStackArgs) {
	affectedBranches := set.New[gitdomain.LocalBranchName]()
	for _, branch := range args.ChangedBranches {
		branchLineage := args.Config.NormalConfig.Lineage.BranchLineageWithoutRoot(branch, args.Config.NormalConfig.PerennialBranches, args.Config.NormalConfig.Order)
		affectedBranches.Add(branchLineage...)
	}
	for _, branch := range affectedBranches.Values() {
		args.Program.Value.Add(&opcodes.ProposalUpdateLineage{
			Branch: branch,
		})
	}
}
