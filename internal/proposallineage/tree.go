package proposallineage

import (
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

type Tree struct {
	Node          *TreeNode
	ProposalCache map[gitdomain.LocalBranchName]Option[forgedomain.Proposal]
}

func NewTree(args ProposalStackLineageArgs) (*Tree, error) {
	tree := &Tree{
		Node: &TreeNode{
			Branch:     "",
			ChildNodes: []*TreeNode{},
			Proposal:   None[forgedomain.Proposal](),
		},
		ProposalCache: map[gitdomain.LocalBranchName]Option[forgedomain.Proposal]{},
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
	descendants := buildAncestorChain(args, self, visited)
	buildDescendantChain(descendants, args, self, visited)
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

func addDescendantNodes(branch gitdomain.LocalBranchName, args ProposalStackLineageArgs, visited map[gitdomain.LocalBranchName]*TreeNode, tree *Tree) {
	if _, ok := visited[branch]; ok {
		return
	}
	parent := args.Lineage.Parent(branch)
	parentBranch, hasParentBranch := parent.Get()
	if !hasParentBranch {
		return
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
		proposal := findProposal(branch, parentBranch, args.Connector)
		branchNode.Proposal = proposal
		tree.ProposalCache[branch] = proposal
	}
	visited[branch] = branchNode
	children := args.Lineage.Children(branch, args.Order)
	for _, child := range children {
		addDescendantNodes(child, args, visited, tree)
	}
}

func buildAncestorChain(
	args ProposalStackLineageArgs,
	tree *Tree,
	visited map[gitdomain.LocalBranchName]*TreeNode,
) gitdomain.LocalBranchNames {
	ancestors := args.Lineage.Ancestors(args.CurrentBranch)
	descendants := gitdomain.LocalBranchNames{args.CurrentBranch}
	previous := tree.Node
	for _, ancestor := range ancestors {
		node := createAncestorNode(ancestor, previous, tree)
		visited[ancestor] = node
		relevantChildren := findRelevantChildren(ancestor, args, ancestors)
		for _, child := range relevantChildren {
			tree.ProposalCache[child.Branch] = child.Proposal
			descendants = append(descendants, child.Branch)
		}
		previous = node
	}
	return descendants
}

func buildDescendantChain(
	descendants gitdomain.LocalBranchNames,
	args ProposalStackLineageArgs,
	tree *Tree,
	visited map[gitdomain.LocalBranchName]*TreeNode,
) {
	for _, descendant := range descendants {
		addDescendantNodes(descendant, args, visited, tree)
	}
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
	sourceBranch gitdomain.LocalBranchName,
	targetBranch gitdomain.LocalBranchName,
	proposalFinder Option[forgedomain.ProposalFinder],
) Option[forgedomain.Proposal] {
	if finder, hasFinder := proposalFinder.Get(); hasFinder {
		proposal, err := finder.FindProposal(sourceBranch, targetBranch)
		if err == nil {
			return proposal
		}
	}
	return None[forgedomain.Proposal]()
}

func findRelevantChildren(
	ancestor gitdomain.LocalBranchName,
	args ProposalStackLineageArgs,
	ancestors gitdomain.LocalBranchNames,
) []ChildWithProposal {
	var result []ChildWithProposal
	for _, child := range args.Lineage.Children(ancestor, args.Order) {
		if shouldIncludeChild(ancestor, child, args.MainAndPerennialBranches, ancestors) {
			proposal := findProposal(child, ancestor, args.Connector)
			result = append(result, ChildWithProposal{
				Branch:   child,
				Proposal: proposal,
			})
		}
	}

	return result
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
