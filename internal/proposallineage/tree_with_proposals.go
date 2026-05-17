package proposallineage

import (
	"github.com/git-town/git-town/v23/internal/forge/forgedomain"
	"github.com/git-town/git-town/v23/internal/git/gitdomain"
	. "github.com/git-town/git-town/v23/pkg/prelude"
)

type TreeNodeWithProposal struct {
	Branch        gitdomain.LocalBranchName
	Children      ForestWithProposals
	LineageParent Option[gitdomain.LocalBranchName]
	Proposal      Option[forgedomain.Proposal]
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

func AddProposalsToTree(tree TreeNode, connector Option[forgedomain.Connector]) TreeNodeWithProposal {
	return addProposalsToTreeHelper(tree, connector)
}

func addProposalsToTreeHelper(tree TreeNode, connector Option[forgedomain.Connector]) TreeNodeWithProposal {
	proposal := loadProposal(tree.Branch, tree.LineageParent, connector)
	children := make(ForestWithProposals, len(tree.Children))
	for i, child := range tree.Children {
		children[i] = addProposalsToTreeHelper(child, connector)
	}
	return TreeNodeWithProposal{
		Branch:        tree.Branch,
		Children:      children,
		LineageParent: tree.LineageParent,
		Proposal:      proposal,
	}
}

func loadProposal(branch gitdomain.LocalBranchName, parentOpt Option[gitdomain.LocalBranchName], connectorOpt Option[forgedomain.Connector]) Option[forgedomain.Proposal] {
	parent, hasParent := parentOpt.Get()
	connector, hasConnector := connectorOpt.Get()
	if !hasParent || !hasConnector {
		return None[forgedomain.Proposal]()
	}
	finder, canFindProposals := connector.(forgedomain.ProposalFinder)
	if !canFindProposals {
		return None[forgedomain.Proposal]()
	}
	proposal, err := finder.FindProposal(branch, parent)
	if err != nil {
		return None[forgedomain.Proposal]()
	}
	return proposal
}
