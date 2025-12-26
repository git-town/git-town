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

func (self TreeNodeWithProposal) BranchOrAncestorHasProposal() bool {
	if self.Proposal.IsSome() {
		return true
	}
	for _, child := range self.Children {
		if child.BranchOrAncestorHasProposal() {
			return true
		}
	}
	return false
}

func AddProposalsToTree(tree TreeNode, proposalFinder Option[forgedomain.ProposalFinder]) TreeNodeWithProposal {
	return addProposalsToTreeHelper(tree, None[gitdomain.LocalBranchName](), proposalFinder)
}

func addProposalsToTreeHelper(tree TreeNode, parent Option[gitdomain.LocalBranchName], connector Option[forgedomain.ProposalFinder]) TreeNodeWithProposal {
	proposal := loadProposal(tree.Branch, parent, connector)
	children := make([]TreeNodeWithProposal, len(tree.Children))
	for i, child := range tree.Children {
		children[i] = addProposalsToTreeHelper(child, Some(tree.Branch), connector)
	}
	return TreeNodeWithProposal{
		Branch:   tree.Branch,
		Children: children,
		Proposal: proposal,
	}
}

func loadProposal(branch gitdomain.LocalBranchName, parentOpt Option[gitdomain.LocalBranchName], connector Option[forgedomain.ProposalFinder]) Option[forgedomain.Proposal] {
	parent, hasParent := parentOpt.Get()
	finder, hasFinder := connector.Get()
	if !hasParent || !hasFinder {
		return None[forgedomain.Proposal]()
	}
	proposal, err := finder.FindProposal(branch, parent)
	if err != nil {
		return None[forgedomain.Proposal]()
	}
	return proposal
}
