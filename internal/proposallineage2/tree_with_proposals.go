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
	return addProposalsToTreeHelper(tree, None[gitdomain.LocalBranchName](), proposalFinder)
}

func addProposalsToTreeHelper(tree TreeNode, parentOpt Option[gitdomain.LocalBranchName], proposalFinder Option[forgedomain.ProposalFinder]) TreeNodeWithProposal {
	parent, hasParent := parentOpt.Get()
	finder, hasFinder := proposalFinder.Get()
	proposal := None[forgedomain.Proposal]()
	if hasParent && hasFinder {
		var err error
		proposal, err = finder.FindProposal(tree.Branch, parent)
		if err != nil {
			proposal = None[forgedomain.Proposal]()
		}
	}
	children := make([]TreeNodeWithProposal, len(tree.Children))
	for i, child := range tree.Children {
		children[i] = addProposalsToTreeHelper(child, Some(tree.Branch), proposalFinder)
	}
	return TreeNodeWithProposal{
		Branch:   tree.Branch,
		Children: children,
		Proposal: proposal,
	}
}
