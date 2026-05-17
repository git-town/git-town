package proposallineage

import (
	"github.com/git-town/git-town/v23/internal/config/configdomain"
	"github.com/git-town/git-town/v23/internal/forge/forgedomain"
	"github.com/git-town/git-town/v23/internal/git/gitdomain"
	. "github.com/git-town/git-town/v23/pkg/prelude"
	"github.com/git-town/git-town/v23/pkg/set"
)

type RenderSectionArgs struct {
	BranchTypes   configdomain.BranchesAndTypes
	Breadcrumb    configdomain.ProposalBreadcrumb
	Connector     Option[forgedomain.Connector]
	CurrentBranch gitdomain.LocalBranchName
	Direction     configdomain.ProposalBreadcrumbDirection
	Excluded      set.Set[configdomain.BranchType]
	Lineage       configdomain.Lineage
	Order         configdomain.Order
}

// RenderSection provides the branch lineage for the given branch in Markdown format, ready to be embedded into a proposal body.
func RenderSection(args RenderSectionArgs) string {
	// step 1: calculate the lineage tree for the given branch
	tree := CalculateTree(args.CurrentBranch, args.Lineage, args.Order, args.BranchTypes)

	// step 2: filter out the excluded branch types
	treeNodes := FilterTree(
		tree,
		args.Excluded,
	)

	branchCount := treeNodes.BranchCount()
	if branchCount == 0 {
		return ""
	}

	if !args.Breadcrumb.DisplayBreadcrumb(branchCount) {
		return ""
	}

	// step 3: add proposals to the tree
	treesWithProposals := treeNodes.AddProposals(args.Connector)

	// step 4: render the tree[s] into Markdown format
	return treesWithProposals.Render(args.CurrentBranch, args.Direction, args.Connector)
}
