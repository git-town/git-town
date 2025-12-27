package sync

import (
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/gohacks/mapstools"
	"github.com/git-town/git-town/v22/internal/proposallineage"
	"github.com/git-town/git-town/v22/internal/vm/opcodes"
	"github.com/git-town/git-town/v22/internal/vm/program"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

type AddStackLineageUpdateOpcodesArgs struct {
	Current                              gitdomain.LocalBranchName
	FullStack                            configdomain.FullStack
	Program                              Mutable[program.Program]
	ProposalStackLineageArgs             proposallineage.ProposalStackLineageArgs
	SkipUpdateForProposalsWithBaseBranch gitdomain.LocalBranchNames
}

// AddStackLineageUpdateOpcodes syncs all given proposals.
// Returns the stack lineage tree if its needed to recall this function.
func AddStackLineageUpdateOpcodes(args AddStackLineageUpdateOpcodesArgs) Option[*proposallineage.Tree] {
	// TODO: there are now multiple places that load and use proposals for branches.
	// To avoid double-loading the same proposal data in one run,
	// extract an object that caches the already known proposals,
	// i.e. which branch has which proposal,
	// and loads missing proposal info on demand.
	tree := proposallineage.NewTree(args.ProposalStackLineageArgs)
	if args.FullStack.Enabled() {
		for branch, proposal := range mapstools.SortedKeyValues(tree.ProposalCache) {
			if args.SkipUpdateForProposalsWithBaseBranch.Contains(branch) {
				continue
			}
			args.Program.Value.Add(&opcodes.ProposalUpdateLineage{
				Branch:         branch,
				Proposal: proposal,
				LineageTree:     MutableSome(tree),
			})
		}
	} else if !args.SkipUpdateForProposalsWithBaseBranch.Contains(args.Current) {
		args.Program.Value.Add(&opcodes.ProposalUpdateLineage{
			Branch:         args.Current,
			Proposal: tree.ProposalCache[args.Current],
			LineageTree:     MutableSome(tree),
		})
	}

	return Some(tree)
}
