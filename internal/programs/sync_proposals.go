package programs

import (
	"github.com/git-town/git-town/v22/internal/config"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/vm/opcodes"
	"github.com/git-town/git-town/v22/internal/vm/program"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

type AddSyncProposalsProgramArgs struct {
	Config          config.ValidatedConfig
	Program         Mutable[program.Program]
	TouchedBranches gitdomain.LocalBranchNames // all branches that the current Git Town command has changed
}

func AddSyncProposalsProgram(args AddSyncProposalsProgramArgs) {
	affectedBranches := args.Config.NormalConfig.Lineage.Clan(args.TouchedBranches, args.Config.MainAndPerennials())
	for _, branch := range affectedBranches {
		args.Program.Value.Add(&opcodes.ProposalUpdateLineage{Branch: branch})
	}
}
