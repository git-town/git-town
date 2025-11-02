package sync

import (
	"fmt"

	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/forge"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/gohacks/mapstools"
	"github.com/git-town/git-town/v22/internal/vm/opcodes"
	"github.com/git-town/git-town/v22/internal/vm/program"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

type AddStackLineageUpdateOpcodesArgs struct {
	Current                              gitdomain.LocalBranchName
	FullStack                            configdomain.FullStack
	Program                              Mutable[program.Program]
	ProposalStackLineageArgs             forge.ProposalStackLineageArgs
	ProposalStackLineageTree             Option[*forge.ProposalStackLineageTree]
	SkipUpdateForProposalsWithBaseBranch gitdomain.LocalBranchNames
}

// AddStackLineageUpdateOpcodes syncs all given proposals.
// Returns the stack lineage tree if its needed to recall this function.
func AddStackLineageUpdateOpcodes(args AddStackLineageUpdateOpcodesArgs) Option[*forge.ProposalStackLineageTree] {
	// TODO: there are now multiple places that load and use proposals for branches.
	// To avoid double-loading the same proposal data in one run,
	// extract an object that caches the already known proposals,
	// i.e. which branch has which proposal,
	// and loads missing proposal info on demand.
	tree, hasTree := args.ProposalStackLineageTree.Get()
	var err error
	if hasTree {
		err = tree.Rebuild(args.ProposalStackLineageArgs)
	} else {
		tree, err = forge.NewProposalStackLineageTree(args.ProposalStackLineageArgs)
	}
	if err != nil {
		fmt.Printf("failed to update proposal stack lineage: %s\n", err.Error())
		return None[*forge.ProposalStackLineageTree]()
	}

	if args.FullStack.Enabled() {
		for branch, proposal := range mapstools.SortedKeyValues(tree.BranchToProposal) {
			if args.SkipUpdateForProposalsWithBaseBranch.Contains(branch) {
				continue
			}
			args.Program.Value.Add(&opcodes.ProposalUpdateLineage{
				Current:         branch,
				CurrentProposal: proposal,
				LineageTree:     MutableSome(tree),
			})
		}
	} else if !args.SkipUpdateForProposalsWithBaseBranch.Contains(args.Current) {
		args.Program.Value.Add(&opcodes.ProposalUpdateLineage{
			Current:         args.Current,
			CurrentProposal: tree.BranchToProposal[args.Current],
			LineageTree:     MutableSome(tree),
		})
	}

	return Some(tree)
}
