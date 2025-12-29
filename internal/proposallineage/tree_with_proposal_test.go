package proposallineage_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/proposallineage"
	"github.com/git-town/git-town/v22/internal/subshell/subshelldomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/shoenig/test/must"
)

// a double that implements the forgedomain.ProposalFinder interface
type testFinder2 struct {
	requests []gitdomain.ProposalTitle
}

func (self *testFinder2) BrowseRepository(_ subshelldomain.Runner) error {
	return nil
}

func (self *testFinder2) CreateProposal(_ forgedomain.CreateProposalArgs) error {
	return nil
}

func (self *testFinder2) DefaultProposalMessage(_ forgedomain.ProposalData) string {
	return ""
}

func (self *testFinder2) FindProposal(source, target gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
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

func TestAddProposalsToTree(t *testing.T) {
	t.Parallel()

	t.Run("all branches have proposals", func(t *testing.T) {
		t.Parallel()
		tree := proposallineage.TreeNode{
			Branch: "main",
			Children: []proposallineage.TreeNode{
				{
					Branch: "feature-a",
					Children: []proposallineage.TreeNode{
						{
							Branch:   "feature-a1",
							Children: []proposallineage.TreeNode{},
						},
					},
				},
				{
					Branch:   "feature-b",
					Children: []proposallineage.TreeNode{},
				},
			},
		}
		var connector forgedomain.Connector = &testFinder2{}
		have := proposallineage.AddProposalsToTree(tree, Some(connector))
		want := proposallineage.TreeNodeWithProposal{
			Branch: "main",
			Children: []proposallineage.TreeNodeWithProposal{
				{
					Branch: "feature-a",
					Children: []proposallineage.TreeNodeWithProposal{
						{
							Branch:   "feature-a1",
							Children: []proposallineage.TreeNodeWithProposal{},
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
					Children: []proposallineage.TreeNodeWithProposal{},
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
		tree := proposallineage.TreeNode{
			Branch: "main",
			Children: []proposallineage.TreeNode{
				{
					Branch: "feature-a",
					Children: []proposallineage.TreeNode{
						{
							Branch:   "feature-a1",
							Children: []proposallineage.TreeNode{},
						},
					},
				},
			},
		}
		var connector forgedomain.Connector = &failingFinder{}
		have := proposallineage.AddProposalsToTree(tree, Some(connector))
		want := proposallineage.TreeNodeWithProposal{
			Branch: "main",
			Children: []proposallineage.TreeNodeWithProposal{
				{
					Branch: "feature-a",
					Children: []proposallineage.TreeNodeWithProposal{
						{
							Branch:   "feature-a1",
							Children: []proposallineage.TreeNodeWithProposal{},
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
		tree := proposallineage.TreeNode{
			Branch: "main",
			Children: []proposallineage.TreeNode{
				{
					Branch: "no-proposal-a",
					Children: []proposallineage.TreeNode{
						{
							Branch:   "feature-a1",
							Children: []proposallineage.TreeNode{},
						},
					},
				},
			},
		}
		var connector forgedomain.Connector = &testFinder2{}
		have := proposallineage.AddProposalsToTree(tree, Some(connector))
		want := proposallineage.TreeNodeWithProposal{
			Branch: "main",
			Children: []proposallineage.TreeNodeWithProposal{
				{
					Branch: "no-proposal-a",
					Children: []proposallineage.TreeNodeWithProposal{
						{
							Branch:   "feature-a1",
							Children: []proposallineage.TreeNodeWithProposal{},
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

	t.Run("multiple children, one has proposal", func(t *testing.T) {
		t.Parallel()
		node := proposallineage.TreeNodeWithProposal{
			Branch: "main",
			Children: []proposallineage.TreeNodeWithProposal{
				{
					Branch:   "feature-a",
					Children: []proposallineage.TreeNodeWithProposal{},
					Proposal: None[forgedomain.Proposal](),
				},
				{
					Branch:   "feature-b",
					Children: []proposallineage.TreeNodeWithProposal{},
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

	t.Run("node has no proposal and children have no proposals", func(t *testing.T) {
		t.Parallel()
		node := proposallineage.TreeNodeWithProposal{
			Branch: "main",
			Children: []proposallineage.TreeNodeWithProposal{
				{
					Branch:   "feature-a",
					Children: []proposallineage.TreeNodeWithProposal{},
					Proposal: None[forgedomain.Proposal](),
				},
				{
					Branch:   "feature-b",
					Children: []proposallineage.TreeNodeWithProposal{},
					Proposal: None[forgedomain.Proposal](),
				},
			},
			Proposal: None[forgedomain.Proposal](),
		}
		must.False(t, node.BranchOrAncestorHasProposal())
	})

	t.Run("node has no proposal and no children", func(t *testing.T) {
		t.Parallel()
		node := proposallineage.TreeNodeWithProposal{
			Branch:   "main",
			Children: []proposallineage.TreeNodeWithProposal{},
			Proposal: None[forgedomain.Proposal](),
		}
		must.False(t, node.BranchOrAncestorHasProposal())
	})

	t.Run("node has no proposal but direct child has proposal", func(t *testing.T) {
		t.Parallel()
		node := proposallineage.TreeNodeWithProposal{
			Branch: "main",
			Children: []proposallineage.TreeNodeWithProposal{
				{
					Branch:   "feature-a",
					Children: []proposallineage.TreeNodeWithProposal{},
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
		node := proposallineage.TreeNodeWithProposal{
			Branch: "main",
			Children: []proposallineage.TreeNodeWithProposal{
				{
					Branch: "feature-a",
					Children: []proposallineage.TreeNodeWithProposal{
						{
							Branch:   "feature-a1",
							Children: []proposallineage.TreeNodeWithProposal{},
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
		node := proposallineage.TreeNodeWithProposal{
			Branch:   "main",
			Children: []proposallineage.TreeNodeWithProposal{},
			Proposal: Some(forgedomain.Proposal{
				Data: forgedomain.ProposalData{
					Title: "test proposal",
				},
			}),
		}
		must.True(t, node.BranchOrAncestorHasProposal())
	})
}
