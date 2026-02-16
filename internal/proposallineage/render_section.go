package proposallineage

import (
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// RenderSection provides the branch lineage for the given branch in Markdown format, ready to be embedded into a proposal body.
func RenderSection(lineage configdomain.Lineage, currentBranch gitdomain.LocalBranchName, order configdomain.Order, breadcrumb configdomain.ProposalBreadcrumb, direction configdomain.ProposalBreadcrumbDirection, connector Option[forgedomain.Connector]) string {
	// step 1: calculate the lineage tree for the given branch
	tree := CalculateTree(currentBranch, lineage, order)

	if !breadcrumb.DisplayBreadcrumb(tree.BranchCount()) {
		return ""
	}

	// step 2: add proposals to the tree
	treeWithProposals := AddProposalsToTree(tree, connector)

	// step 3: render the tree into Markdown format
	return RenderTree(treeWithProposals, currentBranch, direction)
}
