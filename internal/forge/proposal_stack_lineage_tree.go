package forge

import (
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

type ProposalStackLineageTree struct {
	BranchToProposal map[gitdomain.LocalBranchName]Option[forgedomain.Proposal]
	Node             *ProposalStackLineageTreeNode
}

func NewProposalStackLineageTree(args ProposalStackLineageArgs) (*ProposalStackLineageTree, error) {
	tree := &ProposalStackLineageTree{
		BranchToProposal: make(map[gitdomain.LocalBranchName]Option[forgedomain.Proposal]),
		Node:             newProposalStackLineageTreeNode(""),
	}

	err := tree.build(args)
	return tree, err
}

func (self *ProposalStackLineageTree) Rebuild(args ProposalStackLineageArgs) error {
	self.Node = newProposalStackLineageTreeNode("")
	return self.build(args)
}

func (self *ProposalStackLineageTree) build(args ProposalStackLineageArgs) error {
	visited := make(map[gitdomain.LocalBranchName]*ProposalStackLineageTreeNode)

	descendants, err := buildAncestorChain(args, self, visited)
	if err != nil {
		return err
	}

	if err := buildDescendantChain(descendants, args, self, visited); err != nil {
		return err
	}

	if len(self.Node.childNodes) == 0 {
		return nil
	}
	self.Node = self.Node.childNodes[0]
	return nil
}

type ProposalStackLineageTreeNode struct {
	branch     gitdomain.LocalBranchName
	childNodes []*ProposalStackLineageTreeNode
	proposal   Option[forgedomain.Proposal]
}

func addDescendantNodes(branch gitdomain.LocalBranchName, args ProposalStackLineageArgs, visited map[gitdomain.LocalBranchName]*ProposalStackLineageTreeNode, tree *ProposalStackLineageTree) error {
	if _, ok := visited[branch]; ok {
		return nil
	}

	parent := args.Lineage.Parent(branch)
	parentBranch, hasParentBranch := parent.Get()
	if !hasParentBranch {
		return nil
	}
	parentNode := visited[parentBranch]
	branchNode := newProposalStackLineageTreeNode(branch)
	parentNode.childNodes = append(parentNode.childNodes, branchNode)
	if proposal, ok := tree.BranchToProposal[branch]; ok {
		branchNode.proposal = proposal
	} else {
		proposal, err := findProposal(branch, parentBranch, args.Connector)
		if err != nil {
			return err
		}
		branchNode.proposal = proposal
		tree.BranchToProposal[branch] = proposal
	}
	visited[branch] = branchNode

	children := args.Lineage.Children(branch, args.Order)
	for _, child := range children {
		if err := addDescendantNodes(child, args, visited, tree); err != nil {
			return err
		}
	}

	return nil
}

func buildAncestorChain(
	args ProposalStackLineageArgs,
	tree *ProposalStackLineageTree,
	visited map[gitdomain.LocalBranchName]*ProposalStackLineageTreeNode,
) (gitdomain.LocalBranchNames, error) {
	ancestors := args.Lineage.Ancestors(args.CurrentBranch)
	descendants := gitdomain.LocalBranchNames{args.CurrentBranch}
	previous := tree.Node

	for _, ancestor := range ancestors {
		node := createAncestorNode(ancestor, previous, tree)
		visited[ancestor] = node

		relevantChildren, err := findRelevantChildren(ancestor, args, ancestors)
		if err != nil {
			return nil, err
		}

		for _, child := range relevantChildren {
			tree.BranchToProposal[child.branch] = child.proposal
			descendants = append(descendants, child.branch)
		}

		previous = node
	}

	return descendants, nil
}

func buildDescendantChain(
	descendants gitdomain.LocalBranchNames,
	args ProposalStackLineageArgs,
	tree *ProposalStackLineageTree,
	visited map[gitdomain.LocalBranchName]*ProposalStackLineageTreeNode,
) error {
	for _, descendant := range descendants {
		if err := addDescendantNodes(descendant, args, visited, tree); err != nil {
			return err
		}
	}
	return nil
}

type childWithProposal struct {
	branch   gitdomain.LocalBranchName
	proposal Option[forgedomain.Proposal]
}

func createAncestorNode(
	ancestor gitdomain.LocalBranchName,
	parent *ProposalStackLineageTreeNode,
	tree *ProposalStackLineageTree,
) *ProposalStackLineageTreeNode {
	node := newProposalStackLineageTreeNode(ancestor)
	parent.childNodes = append(parent.childNodes, node)
	if proposal, ok := tree.BranchToProposal[ancestor]; ok {
		node.proposal = proposal
	}
	return node
}

func findProposal(
	childBranch gitdomain.LocalBranchName,
	targetBranch gitdomain.LocalBranchName,
	connector Option[forgedomain.ProposalFinder],
) (Option[forgedomain.Proposal], error) {
	if proposalFinder, hasProposalFinder := connector.Get(); hasProposalFinder {
		return proposalFinder.FindProposal(childBranch, targetBranch)
	}
	return None[forgedomain.Proposal](), nil
}

func findRelevantChildren(
	ancestor gitdomain.LocalBranchName,
	args ProposalStackLineageArgs,
	ancestors gitdomain.LocalBranchNames,
) ([]childWithProposal, error) {
	children := args.Lineage.Children(ancestor, args.Order)
	var relevantChildren []childWithProposal

	for _, child := range children {
		if shouldIncludeChild(ancestor, child, args.MainAndPerennialBranches, ancestors) {
			proposal, err := findProposal(child, ancestor, args.Connector)
			if err != nil {
				return nil, err
			}
			relevantChildren = append(relevantChildren, childWithProposal{
				branch:   child,
				proposal: proposal,
			})
		}
	}

	return relevantChildren, nil
}

func newProposalStackLineageTreeNode(branch gitdomain.LocalBranchName) *ProposalStackLineageTreeNode {
	return &ProposalStackLineageTreeNode{
		branch:     branch,
		childNodes: make([]*ProposalStackLineageTreeNode, 0),
		proposal:   None[forgedomain.Proposal](),
	}
}

func shouldIncludeChild(
	ancestor gitdomain.LocalBranchName,
	child gitdomain.LocalBranchName,
	mainAndPerennials gitdomain.LocalBranchNames,
	ancestors gitdomain.LocalBranchNames,
) bool {
	isMainOrPerennial := mainAndPerennials.Contains(ancestor)
	if isMainOrPerennial {
		return ancestors.Contains(child)
	}
	return true
}
