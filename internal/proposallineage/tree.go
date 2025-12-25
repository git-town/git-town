package proposallineage

import (
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

type Tree struct {
	ProposalCache map[gitdomain.LocalBranchName]Option[forgedomain.Proposal]
	Node          *TreeNode
}

func NewTree(args ProposalStackLineageArgs) (*Tree, error) {
	tree := &Tree{
		ProposalCache: map[gitdomain.LocalBranchName]Option[forgedomain.Proposal]{},
		Node: &TreeNode{
			Branch:     "",
			ChildNodes: []*TreeNode{},
			Proposal:   None[forgedomain.Proposal](),
		},
	}
	err := tree.build(args)
	return tree, err
}

func (self *Tree) Rebuild(args ProposalStackLineageArgs) error {
	self.Node = &TreeNode{
		Branch:     "",
		ChildNodes: []*TreeNode{},
		Proposal:   None[forgedomain.Proposal](),
	}
	return self.build(args)
}

func (self *Tree) build(args ProposalStackLineageArgs) error {
	visited := map[gitdomain.LocalBranchName]*TreeNode{}
	descendants, err := buildAncestorChain(args, self, visited)
	if err != nil {
		return err
	}
	if err := buildDescendantChain(descendants, args, self, visited); err != nil {
		return err
	}
	if len(self.Node.ChildNodes) == 0 {
		return nil
	}
	self.Node = self.Node.ChildNodes[0]
	return nil
}

type TreeNode struct {
	Branch     gitdomain.LocalBranchName
	ChildNodes []*TreeNode
	Proposal   Option[forgedomain.Proposal]
}

func addDescendantNodes(branch gitdomain.LocalBranchName, args ProposalStackLineageArgs, visited map[gitdomain.LocalBranchName]*TreeNode, tree *Tree) error {
	if _, ok := visited[branch]; ok {
		return nil
	}
	parent := args.Lineage.Parent(branch)
	parentBranch, hasParentBranch := parent.Get()
	if !hasParentBranch {
		return nil
	}
	parentNode := visited[parentBranch]
	branchNode := &TreeNode{
		Branch:     branch,
		ChildNodes: []*TreeNode{},
		Proposal:   None[forgedomain.Proposal](),
	}
	parentNode.ChildNodes = append(parentNode.ChildNodes, branchNode)
	if proposal, ok := tree.ProposalCache[branch]; ok {
		branchNode.Proposal = proposal
	} else {
		proposal, err := findProposal(branch, parentBranch, args.Connector)
		if err != nil {
			return err
		}
		branchNode.Proposal = proposal
		tree.ProposalCache[branch] = proposal
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
	tree *Tree,
	visited map[gitdomain.LocalBranchName]*TreeNode,
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
			tree.ProposalCache[child.Branch] = child.Proposal
			descendants = append(descendants, child.Branch)
		}
		previous = node
	}
	return descendants, nil
}

func buildDescendantChain(
	descendants gitdomain.LocalBranchNames,
	args ProposalStackLineageArgs,
	tree *Tree,
	visited map[gitdomain.LocalBranchName]*TreeNode,
) error {
	for _, descendant := range descendants {
		if err := addDescendantNodes(descendant, args, visited, tree); err != nil {
			return err
		}
	}
	return nil
}

type ChildWithProposal struct {
	Branch   gitdomain.LocalBranchName
	Proposal Option[forgedomain.Proposal]
}

func createAncestorNode(
	ancestor gitdomain.LocalBranchName,
	parent *TreeNode,
	tree *Tree,
) *TreeNode {
	node := &TreeNode{
		Branch:     ancestor,
		ChildNodes: []*TreeNode{},
		Proposal:   None[forgedomain.Proposal](),
	}
	parent.ChildNodes = append(parent.ChildNodes, node)
	if proposal, ok := tree.ProposalCache[ancestor]; ok {
		node.Proposal = proposal
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
) ([]ChildWithProposal, error) {
	var result []ChildWithProposal
	for _, child := range args.Lineage.Children(ancestor, args.Order) {
		if shouldIncludeChild(ancestor, child, args.MainAndPerennialBranches, ancestors) {
			proposal, err := findProposal(child, ancestor, args.Connector)
			if err != nil {
				return nil, err
			}
			result = append(result, ChildWithProposal{
				Branch:   child,
				Proposal: proposal,
			})
		}
	}

	return result, nil
}

func shouldIncludeChild(
	ancestor gitdomain.LocalBranchName,
	child gitdomain.LocalBranchName,
	mainAndPerennials gitdomain.LocalBranchNames,
	ancestors gitdomain.LocalBranchNames,
) bool {
	if mainAndPerennials.Contains(ancestor) {
		return ancestors.Contains(child)
	}
	return true
}
