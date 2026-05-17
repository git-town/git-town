package proposallineage

import (
	"github.com/git-town/git-town/v23/internal/forge/forgedomain"
	. "github.com/git-town/git-town/v23/pkg/prelude"
)

type TreeNodes []TreeNode

func (self TreeNodes) AddProposals(connector Option[forgedomain.Connector]) TreeNodesWithProposal {
	result := make(TreeNodesWithProposal, 0, len(self))

	for _, tree := range self {
		result = append(result, AddProposalsToTree(tree, connector))
	}
	return result
}

func (self TreeNodes) BranchCount() int {
	var count int
	for _, node := range self {
		count += node.BranchCount()
	}
	return count
}
