package proposallineage2_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/proposallineage2"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/shoenig/test/must"
)

// a double that implements the forgedomain.ProposalFinder interface
type testFinder struct {
	requests []gitdomain.ProposalTitle
}

func (self *testFinder) FindProposal(source, target gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
	if strings.Contains(source.String(), "no-proposal") {
		return None[forgedomain.Proposal](), nil
	}
	title := gitdomain.ProposalTitle(fmt.Sprintf("proposal from %s to %s", source, target))
	self.requests = append(self.requests, title)
	return Some(forgedomain.Proposal{
		Data: forgedomain.ProposalData{
			Title: title,
		},
	}), nil
}

// a Connector double that simulates connection errors
type failingFinder struct{}

func (self *failingFinder) FindProposal(branch, _ gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
	return None[forgedomain.Proposal](), fmt.Errorf("simulated error finding proposal for %s", branch)
}

func TestAddProposalsToTree(t *testing.T) {
	t.Parallel()

	t.Run("all branches have proposals", func(t *testing.T) {
		t.Parallel()
		tree := proposallineage2.TreeNode{
			Branch: "main",
			Children: []proposallineage2.TreeNode{
				{
					Branch: "feature-a",
					Children: []proposallineage2.TreeNode{
						{
							Branch:   "feature-a1",
							Children: []proposallineage2.TreeNode{},
						},
					},
				},
				{
					Branch:   "feature-b",
					Children: []proposallineage2.TreeNode{},
				},
			},
		}
		var connector forgedomain.ProposalFinder = &testFinder{}
		have := proposallineage2.AddProposalsToTree(tree, Some(connector))
		want := proposallineage2.TreeNodeWithProposal{
			Branch: "main",
			Children: []proposallineage2.TreeNodeWithProposal{
				{
					Branch: "feature-a",
					Children: []proposallineage2.TreeNodeWithProposal{
						{
							Branch:   "feature-a1",
							Children: []proposallineage2.TreeNodeWithProposal{},
							Proposal: Some(forgedomain.Proposal{
								Data: forgedomain.ProposalData{
									Title: "proposal from feature-a1 to feature-a",
								},
							}),
						},
					},
					Proposal: Some(forgedomain.Proposal{
						Data: forgedomain.ProposalData{
							Title: "proposal from feature-a to main",
						},
					}),
				},
				{
					Branch:   "feature-b",
					Children: []proposallineage2.TreeNodeWithProposal{},
					Proposal: Some(forgedomain.Proposal{
						Data: forgedomain.ProposalData{
							Title: "proposal from feature-b to main",
						},
					}),
				},
			},
			Proposal: None[forgedomain.Proposal](),
		}
		must.Eq(t, want, have)
	})

	t.Run("connector returns errors", func(t *testing.T) {
		t.Parallel()
		tree := proposallineage2.TreeNode{
			Branch: "main",
			Children: []proposallineage2.TreeNode{
				{
					Branch: "feature-a",
					Children: []proposallineage2.TreeNode{
						{
							Branch:   "feature-a1",
							Children: []proposallineage2.TreeNode{},
						},
					},
				},
			},
		}
		var connector forgedomain.ProposalFinder = &failingFinder{}
		have := proposallineage2.AddProposalsToTree(tree, Some(connector))
		want := proposallineage2.TreeNodeWithProposal{
			Branch: "main",
			Children: []proposallineage2.TreeNodeWithProposal{
				{
					Branch: "feature-a",
					Children: []proposallineage2.TreeNodeWithProposal{
						{
							Branch:   "feature-a1",
							Children: []proposallineage2.TreeNodeWithProposal{},
							Proposal: None[forgedomain.Proposal](),
						},
					},
					Proposal: None[forgedomain.Proposal](),
				},
			},
			Proposal: None[forgedomain.Proposal](),
		}
		must.Eq(t, want, have)
	})

	t.Run("some branches have proposals", func(t *testing.T) {
		t.Parallel()
		tree := proposallineage2.TreeNode{
			Branch: "main",
			Children: []proposallineage2.TreeNode{
				{
					Branch: "no-proposal-a",
					Children: []proposallineage2.TreeNode{
						{
							Branch:   "feature-a1",
							Children: []proposallineage2.TreeNode{},
						},
					},
				},
			},
		}
		var connector forgedomain.ProposalFinder = &testFinder{}
		have := proposallineage2.AddProposalsToTree(tree, Some(connector))
		want := proposallineage2.TreeNodeWithProposal{
			Branch: "main",
			Children: []proposallineage2.TreeNodeWithProposal{
				{
					Branch: "no-proposal-a",
					Children: []proposallineage2.TreeNodeWithProposal{
						{
							Branch:   "feature-a1",
							Children: []proposallineage2.TreeNodeWithProposal{},
							Proposal: Some(forgedomain.Proposal{
								Data: forgedomain.ProposalData{
									Title: "proposal from feature-a1 to no-proposal-a",
								},
							}),
						},
					},
					Proposal: None[forgedomain.Proposal](),
				},
			},
			Proposal: None[forgedomain.Proposal](),
		}
		must.Eq(t, want, have)
	})
}

func TestBranchOrAncestorHasProposal(t *testing.T) {
	t.Parallel()

	t.Run("node has no proposal and children have no proposals", func(t *testing.T) {
		t.Parallel()
		node := proposallineage2.TreeNodeWithProposal{
			Branch: "main",
			Children: []proposallineage2.TreeNodeWithProposal{
				{
					Branch:   "feature-a",
					Children: []proposallineage2.TreeNodeWithProposal{},
					Proposal: None[forgedomain.Proposal](),
				},
				{
					Branch:   "feature-b",
					Children: []proposallineage2.TreeNodeWithProposal{},
					Proposal: None[forgedomain.Proposal](),
				},
			},
			Proposal: None[forgedomain.Proposal](),
		}
		must.False(t, node.BranchOrAncestorHasProposal())
	})

	t.Run("node has no proposal and no children", func(t *testing.T) {
		t.Parallel()
		node := proposallineage2.TreeNodeWithProposal{
			Branch:   "main",
			Children: []proposallineage2.TreeNodeWithProposal{},
			Proposal: None[forgedomain.Proposal](),
		}
		must.False(t, node.BranchOrAncestorHasProposal())
	})

	t.Run("multiple children, one has proposal", func(t *testing.T) {
		t.Parallel()
		node := proposallineage2.TreeNodeWithProposal{
			Branch: "main",
			Children: []proposallineage2.TreeNodeWithProposal{
				{
					Branch:   "feature-a",
					Children: []proposallineage2.TreeNodeWithProposal{},
					Proposal: None[forgedomain.Proposal](),
				},
				{
					Branch:   "feature-b",
					Children: []proposallineage2.TreeNodeWithProposal{},
					Proposal: Some(forgedomain.Proposal{
						Data: forgedomain.ProposalData{
							Title: "proposal for feature-b",
						},
					}),
				},
			},
			Proposal: None[forgedomain.Proposal](),
		}
		must.True(t, node.BranchOrAncestorHasProposal())
	})

	t.Run("node has no proposal but direct child has proposal", func(t *testing.T) {
		t.Parallel()
		node := proposallineage2.TreeNodeWithProposal{
			Branch: "main",
			Children: []proposallineage2.TreeNodeWithProposal{
				{
					Branch:   "feature-a",
					Children: []proposallineage2.TreeNodeWithProposal{},
					Proposal: Some(forgedomain.Proposal{
						Data: forgedomain.ProposalData{
							Title: "proposal for feature-a",
						},
					}),
				},
			},
			Proposal: None[forgedomain.Proposal](),
		}
		must.True(t, node.BranchOrAncestorHasProposal())
	})

	t.Run("node has no proposal but grandchild has proposal", func(t *testing.T) {
		t.Parallel()
		node := proposallineage2.TreeNodeWithProposal{
			Branch: "main",
			Children: []proposallineage2.TreeNodeWithProposal{
				{
					Branch: "feature-a",
					Children: []proposallineage2.TreeNodeWithProposal{
						{
							Branch:   "feature-a1",
							Children: []proposallineage2.TreeNodeWithProposal{},
							Proposal: Some(forgedomain.Proposal{
								Data: forgedomain.ProposalData{
									Title: "proposal for feature-a1",
								},
							}),
						},
					},
					Proposal: None[forgedomain.Proposal](),
				},
			},
			Proposal: None[forgedomain.Proposal](),
		}
		must.True(t, node.BranchOrAncestorHasProposal())
	})

	t.Run("node has proposal", func(t *testing.T) {
		t.Parallel()
		node := proposallineage2.TreeNodeWithProposal{
			Branch:   "main",
			Children: []proposallineage2.TreeNodeWithProposal{},
			Proposal: Some(forgedomain.Proposal{
				Data: forgedomain.ProposalData{
					Title: "test proposal",
				},
			}),
		}
		must.True(t, node.BranchOrAncestorHasProposal())
	})
}
