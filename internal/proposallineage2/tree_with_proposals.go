package proposallineage2

import (
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

type TreeNodeWithProposal struct {
	Branch   gitdomain.LocalBranchName
	Children []TreeNodeWithProposal
	Proposal Option[forgedomain.Proposal]
}

func AddProposalsToTree(tree TreeNode, proposalFinder Option[forgedomain.ProposalFinder]) TreeNodeWithProposal {
	return TreeNodeWithProposal{
		Branch:   tree.Branch,
		Children: []TreeNodeWithProposal{},
		Proposal: Option[forgedomain.Proposal]{},
	}
}
