package sync

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/forge"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/vm/opcodes"
	"github.com/git-town/git-town/v21/internal/vm/program"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

type BranchProposalsProgramArgs struct {
	Current                  gitdomain.LocalBranchName
	FullStack                configdomain.FullStack
	Program                  Mutable[program.Program]
	ProposalStackLineageArgs forge.ProposalStackLineageArgs
}

// BranchProposalsProgram syncs all given proposals.
func BranchProposalsProgram(args BranchProposalsProgramArgs) {
	// TODO: there are now multiple places that load and use proposals for branches.
	// To avoid double-loading the same proposal data in one run,
	// extract an object that caches the already known proposals,
	// i.e. which branch has which proposal,
	// and loads missing proposal info on demand.
	tree, err := forge.NewProposalStackLineageTree(args.ProposalStackLineageArgs)
	if err != nil {
		fmt.Printf("failed to update proposal stack lineage: %s\n", err.Error())
		return
	}

	if args.FullStack.Enabled() {
		for branch, proposal := range tree.BranchToProposal {
			args.Program.Value.Add(&opcodes.ProposalUpdateLineage{
				Current:         branch,
				CurrentProposal: proposal,
				LineageTree:     MutableSome(tree),
			})
		}
	} else {
		args.Program.Value.Add(&opcodes.ProposalUpdateLineage{
			Current:         args.Current,
			CurrentProposal: tree.BranchToProposal[args.Current],
			LineageTree:     MutableSome(tree),
		})
	}
}
