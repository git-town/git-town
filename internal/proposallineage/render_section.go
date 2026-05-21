package proposallineage

import (
	"github.com/git-town/git-town/v23/internal/config/configdomain"
	"github.com/git-town/git-town/v23/internal/forge/forgedomain"
	"github.com/git-town/git-town/v23/internal/git/gitdomain"
	. "github.com/git-town/git-town/v23/pkg/prelude"
)

type RenderSectionArgs struct {
	Breadcrumb    configdomain.ProposalBreadcrumb
	Connector     Option[forgedomain.Connector]
	CurrentBranch gitdomain.LocalBranchName
	Direction     configdomain.ProposalBreadcrumbDirection
	Lineage       configdomain.Lineage
	Order         configdomain.Order
}

// RenderSection provides the branch lineage for the given branch in Markdown format, ready to be embedded into a proposal body.
func RenderSection(args RenderSectionArgs) string {
	// step 1: calculate the lineage tree for the given branch
	tree := CalculateTree(args.CurrentBranch, args.Lineage, args.Order)

	if !args.Breadcrumb.DisplayBreadcrumb(tree.BranchCount()) {
		return ""
	}

	// step 2: add proposals to the tree
	treeWithProposals := AddProposalsToTree(tree, args.Connector)

	// step 3: render the tree into Markdown format
	return RenderTree(treeWithProposals, args.CurrentBranch, args.Direction, args.Connector)
}
