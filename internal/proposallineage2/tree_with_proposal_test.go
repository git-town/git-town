package proposallineage2

import (
	"fmt"
	"strings"
	"testing"

	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
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
		tree := TreeNode{
			Branch: "main",
			Children: []TreeNode{
				{
					Branch: "feature-a",
					Children: []TreeNode{
						{
							Branch:   "feature-a1",
							Children: []TreeNode{},
						},
					},
				},
				{
					Branch:   "feature-b",
					Children: []TreeNode{},
				},
			},
		}
		var connector forgedomain.ProposalFinder = &testFinder{}
		have := AddProposalsToTree(tree, Some(connector))
		want := TreeNodeWithProposal{
			Branch: "main",
			Children: []TreeNodeWithProposal{
				{
					Branch: "feature-a",
					Children: []TreeNodeWithProposal{
						{
							Branch:   "feature-a1",
							Children: []TreeNodeWithProposal{},
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
					Children: []TreeNodeWithProposal{},
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
		tree := TreeNode{
			Branch: "main",
			Children: []TreeNode{
				{
					Branch: "feature-a",
					Children: []TreeNode{
						{
							Branch:   "feature-a1",
							Children: []TreeNode{},
						},
					},
				},
			},
		}
		var connector forgedomain.ProposalFinder = &failingFinder{}
		have := AddProposalsToTree(tree, Some(connector))
		want := TreeNodeWithProposal{
			Branch: "main",
			Children: []TreeNodeWithProposal{
				{
					Branch: "no-proposal-a",
					Children: []TreeNodeWithProposal{
						{
							Branch:   "feature-a1",
							Children: []TreeNodeWithProposal{},
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
		tree := TreeNode{
			Branch: "main",
			Children: []TreeNode{
				{
					Branch: "no-proposal-a",
					Children: []TreeNode{
						{
							Branch:   "feature-a1",
							Children: []TreeNode{},
						},
					},
				},
			},
		}
		var connector forgedomain.ProposalFinder = &testFinder{}
		have := AddProposalsToTree(tree, Some(connector))
		want := TreeNodeWithProposal{
			Branch: "main",
			Children: []TreeNodeWithProposal{
				{
					Branch: "no-proposal-a",
					Children: []TreeNodeWithProposal{
						{
							Branch:   "feature-a1",
							Children: []TreeNodeWithProposal{},
							Proposal: Some(forgedomain.Proposal{
								Data: forgedomain.ProposalData{
									Title: "proposal from feature-a1 to feature-a",
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
