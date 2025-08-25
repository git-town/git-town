package sync

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/forge"
	"github.com/git-town/git-town/v21/internal/vm/opcodes"
	"github.com/git-town/git-town/v21/internal/vm/program"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

type BranchProposalsProgramArgs struct {
	Program                  Mutable[program.Program]
	ProposalStackLineageArgs forge.ProposalStackLineageArgs
}

// BranchProposalsProgram syncs all given proposals.
func BranchProposalsProgram(branchesToSync configdomain.BranchesToSync, args BranchProposalsProgramArgs) {
	tree, err := forge.NewProposalStackLineageTree(args.ProposalStackLineageArgs)
	if err != nil {
		fmt.Printf("failed to update proposal stack lineage: %s\n", err.Error())
		return
	}

	for _, branch := range branchesToSync {
		// TODO: there are now multiple places that load and use proposals for branches.
		// To avoid double-loading the same proposal data in one run,
		// extract an object that caches the already known proposals,
		// i.e. which branch has which proposal,
		// and loads missing proposal info on demand.
		proposal, ok := tree.BranchToProposal[branch.BranchInfo.LocalBranchName()]
		if !ok {
			continue
		}
		args.Program.Value.Add(&opcodes.ProposalUpdateLineage{
			Current:         branch.BranchInfo.LocalBranchName(),
			CurrentProposal: proposal,
			LineageTree:     MutableSome(tree),
		})
	}
}
