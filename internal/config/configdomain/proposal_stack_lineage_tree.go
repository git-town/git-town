package configdomain

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

func NewProposalStackLineageTree(args ProposalStackLineageArgs) (*ProposalStackLineageTree, error) {
	tree := ProposalStackLineageTree{
		BranchToProposal: make(map[gitdomain.LocalBranchName]Option[forgedomain.Proposal]),
		Node:             newProposalStackLineageTreeNode(""),
	}

	// Start of DFS
	visited := make(map[gitdomain.LocalBranchName]*ProposalStackLineageTreeNode)
	ancestors := args.Lineage.Ancestors(args.CurrentBranch)
	descendantsOfAncestor := gitdomain.NewLocalBranchNames(args.CurrentBranch.String())
	previous := tree.Node
	var ancestorTreeNode *ProposalStackLineageTreeNode

	// climb up the lineage chain and find all corresponding proposals
	for _, ancestor := range ancestors {
		ancestorTreeNode = newProposalStackLineageTreeNode(ancestor)
		previous.childNodes = append(previous.childNodes, ancestorTreeNode)
		ancestorTreeNode.depth = previous.depth + 1
		if proposal, ok := tree.BranchToProposal[ancestor]; ok {
			ancestorTreeNode.proposal = proposal
		}
		visited[ancestor] = ancestorTreeNode

		children := args.Lineage.Children(ancestor)
		for _, child := range children {
			if !args.MainAndPerennialBranches.Contains(child) && !ancestors.Contains(child) {
				childBranchProposal, err := findProposal(child, ancestor, args.Connector)
				if err != nil {
					return nil, err
				}
				tree.BranchToProposal[child] = childBranchProposal
				descendantsOfAncestor = append(descendantsOfAncestor, child)
			}
		}
		previous = ancestorTreeNode
	}

	// Next get all descendants
	for _, descendant := range descendantsOfAncestor {
		if err := addDescendantNodes(descendant, args, visited, &tree); err != nil {
			return nil, err
		}
	}

	tree.Node = tree.Node.childNodes[0]
	return &tree, nil
}

type ProposalStackLineageTree struct {
	BranchToProposal map[gitdomain.LocalBranchName]Option[forgedomain.Proposal]
	Node             *ProposalStackLineageTreeNode
}

type ProposalStackLineageTreeNode struct {
	branch     gitdomain.LocalBranchName
	childNodes []*ProposalStackLineageTreeNode
	depth      int
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
	branchNode.depth = parentNode.depth + 1
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

	children := args.Lineage.Children(branch)
	for _, child := range children {
		if err := addDescendantNodes(child, args, visited, tree); err != nil {
			return err
		}
	}

	return nil
}

func findProposal(
	childBranch gitdomain.LocalBranchName,
	targetBranch gitdomain.LocalBranchName,
	connector forgedomain.Connector,
) (Option[forgedomain.Proposal], error) {
	findProposalFn, hasFindProposalFn := connector.FindProposalFn().Get()
	if !hasFindProposalFn {
		return None[forgedomain.Proposal](), nil
	}

	proposal, err := findProposalFn(childBranch, targetBranch)
	if err != nil {
		return None[forgedomain.Proposal](), fmt.Errorf("failed to find proposal for branch %s: %w", childBranch, err)
	}

	proposalData, hasProposal := proposal.Get()
	if !hasProposal {
		return None[forgedomain.Proposal](), nil
	}

	return Some(proposalData), nil
}

func newProposalStackLineageTreeNode(branch gitdomain.LocalBranchName) *ProposalStackLineageTreeNode {
	return &ProposalStackLineageTreeNode{
		branch:     branch,
		childNodes: make([]*ProposalStackLineageTreeNode, 0),
		depth:      -1,
		proposal:   None[forgedomain.Proposal](),
	}
}
