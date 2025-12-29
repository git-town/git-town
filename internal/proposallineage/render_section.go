package proposallineage

import (
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// RenderSection provides the branch lineage for the given branch in Markdown format, ready to be embedded into a proposal body.
func RenderSection(lineage configdomain.Lineage, currentBranch gitdomain.LocalBranchName, order configdomain.Order, connector Option[forgedomain.Connector]) string {
	// step 1: calculate the structure of the lineage tree to display
	tree := CalculateTree(currentBranch, lineage, order)

	// step 2: add proposals to the tree
	treeWithProposals := AddProposalsToTree(tree, connector)

	// step 3: render the tree into Markdown
	return RenderTree(treeWithProposals, currentBranch)
}
