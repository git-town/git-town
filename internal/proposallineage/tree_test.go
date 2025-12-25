package proposallineage_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/proposallineage"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/shoenig/test/must"
)

// a Connector double that implements the forgedomain.ProposalFinder interface
type testConnector struct {
	requests []gitdomain.ProposalTitle
}

func (self *testConnector) FindProposal(source, target gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
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
type failingConnector struct{}

func (self *failingConnector) FindProposal(branch, _ gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
	return None[forgedomain.Proposal](), fmt.Errorf("mock error finding proposal for %s", branch)
}

func TestNewTree(t *testing.T) {
	t.Parallel()

	t.Run("on feature branch with multiple children", func(t *testing.T) {
		t.Parallel()
		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			"feature-a": "main",
			"feature-b": "feature-a",
			"feature-c": "feature-a",
		})
		connector := testConnector{}
		var proposalFinder forgedomain.ProposalFinder = &connector
		have, err := proposallineage.NewTree(proposallineage.ProposalStackLineageArgs{
			Connector:                Some(proposalFinder),
			CurrentBranch:            "feature-a",
			Lineage:                  lineage,
			MainAndPerennialBranches: gitdomain.LocalBranchNames{"main"},
		})
		want := &proposallineage.Tree{
			BranchToProposal: map[gitdomain.LocalBranchName]Option[forgedomain.Proposal]{
				"feature-a": Some(forgedomain.Proposal{
					Data: forgedomain.ProposalData{
						Title: "proposal from feature-a to main",
					},
				}),
				"feature-b": Some(forgedomain.Proposal{
					Data: forgedomain.ProposalData{
						Title: "proposal from feature-b to feature-a",
					},
				}),
				"feature-c": Some(forgedomain.Proposal{
					Data: forgedomain.ProposalData{
						Title: "proposal from feature-c to feature-a",
					},
				}),
			},
			Node: &proposallineage.TreeNode{
				Branch: "main",
				ChildNodes: []*proposallineage.TreeNode{
					{
						Branch: "feature-a",
						ChildNodes: []*proposallineage.TreeNode{
							{
								Branch:     "feature-b",
								ChildNodes: []*proposallineage.TreeNode{},
								Proposal: Some(forgedomain.Proposal{
									Data: forgedomain.ProposalData{
										Title: "proposal from feature-b to feature-a",
									},
								}),
							},
							{
								Branch:     "feature-c",
								ChildNodes: []*proposallineage.TreeNode{},
								Proposal: Some(forgedomain.Proposal{
									Data: forgedomain.ProposalData{
										Title: "proposal from feature-c to feature-a",
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
				},
				Proposal: None[forgedomain.Proposal](),
			},
		}
		must.NoError(t, err)
		must.Eq(t, want, have)
		wantRequests := []gitdomain.ProposalTitle{
			"proposal from feature-a to main",
			"proposal from feature-b to feature-a",
			"proposal from feature-c to feature-a",
		}
		must.Eq(t, wantRequests, connector.requests)
	})

	t.Run("on a feature branch in the middle of a long lineage", func(t *testing.T) {
		t.Parallel()
		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			"feature-a": "main",
			"feature-b": "feature-a",
			"feature-c": "feature-b",
			"feature-d": "feature-c",
		})
		var connector forgedomain.ProposalFinder = &testConnector{}
		have, err := proposallineage.NewTree(proposallineage.ProposalStackLineageArgs{
			Connector:                Some(connector),
			CurrentBranch:            "feature-a",
			Lineage:                  lineage,
			MainAndPerennialBranches: gitdomain.LocalBranchNames{"main"},
		})
		want := &proposallineage.Tree{
			BranchToProposal: map[gitdomain.LocalBranchName]Option[forgedomain.Proposal]{
				"feature-a": Some(forgedomain.Proposal{
					Data: forgedomain.ProposalData{
						Title: "proposal from feature-a to main",
					},
				}),
				"feature-b": Some(forgedomain.Proposal{
					Data: forgedomain.ProposalData{
						Title: "proposal from feature-b to feature-a",
					},
				}),
				"feature-c": Some(forgedomain.Proposal{
					Data: forgedomain.ProposalData{
						Title: "proposal from feature-c to feature-b",
					},
				}),
				"feature-d": Some(forgedomain.Proposal{
					Data: forgedomain.ProposalData{
						Title: "proposal from feature-d to feature-c",
					},
				}),
			},
			Node: &proposallineage.TreeNode{
				Branch: "main",
				ChildNodes: []*proposallineage.TreeNode{
					{
						Branch: "feature-a",
						ChildNodes: []*proposallineage.TreeNode{
							{
								Branch: "feature-b",
								ChildNodes: []*proposallineage.TreeNode{
									{
										Branch: "feature-c",
										ChildNodes: []*proposallineage.TreeNode{
											{
												Branch:     "feature-d",
												ChildNodes: []*proposallineage.TreeNode{},
												Proposal: Some(forgedomain.Proposal{
													Data: forgedomain.ProposalData{
														Title: "proposal from feature-d to feature-c",
													},
												}),
											},
										},
										Proposal: Some(forgedomain.Proposal{
											Data: forgedomain.ProposalData{
												Title: "proposal from feature-c to feature-b",
											},
										}),
									},
								},
								Proposal: Some(forgedomain.Proposal{
									Data: forgedomain.ProposalData{
										Title: "proposal from feature-b to feature-a",
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
				},
				Proposal: None[forgedomain.Proposal](),
			},
		}
		must.NoError(t, err)
		must.Eq(t, want, have)
	})

	t.Run("on a leaf branch with siblings", func(t *testing.T) {
		t.Parallel()
		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			"feature-a": "main",
			"feature-b": "feature-a",
			"feature-c": "feature-a",
		})
		var connector forgedomain.ProposalFinder = &testConnector{}
		have, err := proposallineage.NewTree(proposallineage.ProposalStackLineageArgs{
			Connector:                Some(connector),
			CurrentBranch:            "feature-c",
			Lineage:                  lineage,
			MainAndPerennialBranches: gitdomain.LocalBranchNames{"main"},
		})
		want := &proposallineage.Tree{
			BranchToProposal: map[gitdomain.LocalBranchName]Option[forgedomain.Proposal]{
				"feature-a": Some(forgedomain.Proposal{
					Data: forgedomain.ProposalData{
						Title: "proposal from feature-a to main",
					},
				}),
				"feature-b": Some(forgedomain.Proposal{
					Data: forgedomain.ProposalData{
						Title: "proposal from feature-b to feature-a",
					},
				}),
				"feature-c": Some(forgedomain.Proposal{
					Data: forgedomain.ProposalData{
						Title: "proposal from feature-c to feature-a",
					},
				}),
			},
			Node: &proposallineage.TreeNode{
				Branch: "main",
				ChildNodes: []*proposallineage.TreeNode{
					{
						Branch: "feature-a",
						ChildNodes: []*proposallineage.TreeNode{
							{
								Branch:     "feature-c",
								ChildNodes: []*proposallineage.TreeNode{},
								Proposal: Some(forgedomain.Proposal{
									Data: forgedomain.ProposalData{
										Title: "proposal from feature-c to feature-a",
									},
								}),
							},
							{
								Branch:     "feature-b",
								ChildNodes: []*proposallineage.TreeNode{},
								Proposal: Some(forgedomain.Proposal{
									Data: forgedomain.ProposalData{
										Title: "proposal from feature-b to feature-a",
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
				},
				Proposal: None[forgedomain.Proposal](),
			},
		}
		must.NoError(t, err)
		must.Eq(t, want, have)
	})

	t.Run("on the perennial branch at the root", func(t *testing.T) {
		t.Parallel()
		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			"feature-a": "main",
			"feature-b": "feature-a",
		})
		var connector forgedomain.ProposalFinder = &testConnector{}
		have, err := proposallineage.NewTree(proposallineage.ProposalStackLineageArgs{
			Connector:                Some(connector),
			CurrentBranch:            "main",
			Lineage:                  lineage,
			MainAndPerennialBranches: gitdomain.LocalBranchNames{"main"},
		})
		want := &proposallineage.Tree{
			BranchToProposal: map[gitdomain.LocalBranchName]Option[forgedomain.Proposal]{},
			Node: &proposallineage.TreeNode{
				Branch:     "",
				ChildNodes: []*proposallineage.TreeNode{},
				Proposal:   None[forgedomain.Proposal](),
			},
		}
		must.NoError(t, err)
		must.Eq(t, want, have)
	})

	t.Run("simple feature branch", func(t *testing.T) {
		t.Parallel()
		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			"feature-a": "main",
		})
		var connector forgedomain.ProposalFinder = &testConnector{}
		have, err := proposallineage.NewTree(proposallineage.ProposalStackLineageArgs{
			Connector:                Some(connector),
			CurrentBranch:            "feature-a",
			Lineage:                  lineage,
			MainAndPerennialBranches: gitdomain.LocalBranchNames{"main"},
		})
		want := &proposallineage.Tree{
			BranchToProposal: map[gitdomain.LocalBranchName]Option[forgedomain.Proposal]{
				"feature-a": Some(forgedomain.Proposal{
					Data: forgedomain.ProposalData{
						Title: "proposal from feature-a to main",
					},
				}),
			},
			Node: &proposallineage.TreeNode{
				Branch: "main",
				ChildNodes: []*proposallineage.TreeNode{
					{
						Branch:     "feature-a",
						ChildNodes: []*proposallineage.TreeNode{},
						Proposal: Some(forgedomain.Proposal{
							Data: forgedomain.ProposalData{
								Title: "proposal from feature-a to main",
							},
						}),
					},
				},
				Proposal: None[forgedomain.Proposal](),
			},
		}
		must.NoError(t, err)
		must.Eq(t, want, have)
	})

	t.Run("some branches have no proposal", func(t *testing.T) {
		t.Parallel()
		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			"feature-a":     "main",
			"no-proposal-b": "feature-a",
		})
		var connector forgedomain.ProposalFinder = &testConnector{}
		have, err := proposallineage.NewTree(proposallineage.ProposalStackLineageArgs{
			Connector:                Some(connector),
			CurrentBranch:            "feature-a",
			Lineage:                  lineage,
			MainAndPerennialBranches: gitdomain.LocalBranchNames{"main"},
		})
		want := &proposallineage.Tree{
			BranchToProposal: map[gitdomain.LocalBranchName]Option[forgedomain.Proposal]{
				"feature-a": Some(forgedomain.Proposal{
					Data: forgedomain.ProposalData{
						Title: "proposal from feature-a to main",
					},
				}),
				"no-proposal-b": None[forgedomain.Proposal](),
			},
			Node: &proposallineage.TreeNode{
				Branch: "main",
				ChildNodes: []*proposallineage.TreeNode{
					{
						Branch: "feature-a",
						ChildNodes: []*proposallineage.TreeNode{
							{
								Branch:     "no-proposal-b",
								ChildNodes: []*proposallineage.TreeNode{},
								Proposal:   None[forgedomain.Proposal](),
							},
						},
						Proposal: Some(forgedomain.Proposal{
							Data: forgedomain.ProposalData{
								Title: "proposal from feature-a to main",
							},
						}),
					},
				},
				Proposal: None[forgedomain.Proposal](),
			},
		}
		must.NoError(t, err)
		must.Eq(t, want, have)
	})

	t.Run("no connector", func(t *testing.T) {
		t.Parallel()
		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			"feature-a": "main",
			"feature-b": "feature-a",
		})
		have, err := proposallineage.NewTree(proposallineage.ProposalStackLineageArgs{
			Connector:                None[forgedomain.ProposalFinder](),
			CurrentBranch:            "feature-a",
			Lineage:                  lineage,
			MainAndPerennialBranches: gitdomain.LocalBranchNames{"main"},
		})
		want := &proposallineage.Tree{
			BranchToProposal: map[gitdomain.LocalBranchName]Option[forgedomain.Proposal]{
				"feature-a": None[forgedomain.Proposal](),
				"feature-b": None[forgedomain.Proposal](),
			},
			Node: &proposallineage.TreeNode{
				Branch: "main",
				ChildNodes: []*proposallineage.TreeNode{
					{
						Branch: "feature-a",
						ChildNodes: []*proposallineage.TreeNode{
							{
								Branch:     "feature-b",
								ChildNodes: []*proposallineage.TreeNode{},
								Proposal:   None[forgedomain.Proposal](),
							},
						},
						Proposal: None[forgedomain.Proposal](),
					},
				},
				Proposal: None[forgedomain.Proposal](),
			},
		}
		must.NoError(t, err)
		must.Eq(t, want, have)
	})

	t.Run("connector returns errors", func(t *testing.T) {
		t.Parallel()
		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			"feature-a": "main",
			"feature-b": "feature-a",
		})
		var connector forgedomain.ProposalFinder = &failingConnector{}
		have, err := proposallineage.NewTree(proposallineage.ProposalStackLineageArgs{
			Connector:                Some(connector),
			CurrentBranch:            "feature-a",
			Lineage:                  lineage,
			MainAndPerennialBranches: gitdomain.LocalBranchNames{"main"},
		})
		want := &proposallineage.Tree{
			BranchToProposal: map[gitdomain.LocalBranchName]Option[forgedomain.Proposal]{},
			Node: &proposallineage.TreeNode{
				Branch: "", // TODO: this shouldn't be empty
				ChildNodes: []*proposallineage.TreeNode{
					{
						Branch: "main",
						ChildNodes: []*proposallineage.TreeNode{
							{
								Branch:     "feature-a",
								ChildNodes: []*proposallineage.TreeNode{},
								Proposal:   None[forgedomain.Proposal](),
							},
						},
						Proposal: None[forgedomain.Proposal](),
					},
				},
				Proposal: None[forgedomain.Proposal](),
			},
		}
		must.Error(t, err) // TODO: should it ignore errors and create the lineage without proposals?
		must.Eq(t, want, have)
	})
}

func TestTreeRebuild(t *testing.T) {
	t.Parallel()

	t.Run("builds a new tree using the cached proposals", func(t *testing.T) {
		t.Parallel()
		// build tree for lineage 1
		lineage1 := configdomain.NewLineageWith(configdomain.LineageData{
			"feature-a": "main",
			"feature-b": "feature-a",
		})
		connector := testConnector{}
		var proposalFinder forgedomain.ProposalFinder = &connector
		tree, err := proposallineage.NewTree(proposallineage.ProposalStackLineageArgs{
			Connector:                Some(proposalFinder),
			CurrentBranch:            "feature-a",
			Lineage:                  lineage1,
			MainAndPerennialBranches: gitdomain.LocalBranchNames{"main"},
		})
		must.NoError(t, err)
		wantRequests := []gitdomain.ProposalTitle{
			"proposal from feature-a to main",
			"proposal from feature-b to feature-a",
		}
		must.Eq(t, wantRequests, connector.requests)
		// build tree for lineage 2
		lineage2 := configdomain.NewLineageWith(configdomain.LineageData{
			"feature-a": "main",
			"feature-b": "feature-a",
			"feature-c": "feature-b",
		})
		connector = testConnector{}
		proposalFinder = &connector
		err = tree.Rebuild(proposallineage.ProposalStackLineageArgs{
			Connector:                Some(proposalFinder),
			CurrentBranch:            "feature-b",
			Lineage:                  lineage2,
			MainAndPerennialBranches: gitdomain.LocalBranchNames{"main"},
		})
		must.NoError(t, err)
		must.NotNil(t, tree.Node)
		// TODO: Make it not look up the proposals for feature-a and feature-b again. They should be cached internally.
		wantRequests = []gitdomain.ProposalTitle{
			"proposal from feature-a to main",
			"proposal from feature-b to feature-a",
			"proposal from feature-c to feature-b",
		}
		must.Eq(t, wantRequests, connector.requests)
	})

	t.Run("error during rebuild", func(t *testing.T) {
		t.Parallel()
		// build tree for lineage 1
		lineage1 := configdomain.NewLineageWith(configdomain.LineageData{
			"feature-a": "main",
			"feature-b": "feature-a",
		})
		connector := testConnector{}
		var proposalFinder forgedomain.ProposalFinder = &connector
		tree, err := proposallineage.NewTree(proposallineage.ProposalStackLineageArgs{
			Connector:                Some(proposalFinder),
			CurrentBranch:            "feature-a",
			Lineage:                  lineage1,
			MainAndPerennialBranches: gitdomain.LocalBranchNames{"main"},
		})
		must.NoError(t, err)
		wantRequests := []gitdomain.ProposalTitle{
			"proposal from feature-a to main",
			"proposal from feature-b to feature-a",
		}
		must.Eq(t, wantRequests, connector.requests)
		// build tree for lineage 2
		lineage2 := configdomain.NewLineageWith(configdomain.LineageData{
			"feature-a": "main",
			"feature-b": "feature-a",
			"feature-c": "feature-b",
		})
		var errorConnector forgedomain.ProposalFinder = &failingConnector{}
		err = tree.Rebuild(proposallineage.ProposalStackLineageArgs{
			Connector:                Some(errorConnector),
			CurrentBranch:            "feature-b",
			Lineage:                  lineage2,
			MainAndPerennialBranches: gitdomain.LocalBranchNames{"main"},
		})
		must.Error(t, err)
	})
}

func TestTreeWithComplexLineages(t *testing.T) {
	t.Parallel()

	t.Run("handles multiple children with different order settings", func(t *testing.T) {
		t.Parallel()
		mainBranch := gitdomain.NewLocalBranchName("main")
		parent := gitdomain.NewLocalBranchName("parent")
		childA := gitdomain.NewLocalBranchName("child-a")
		childB := gitdomain.NewLocalBranchName("child-b")
		childC := gitdomain.NewLocalBranchName("child-c")

		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			childA: parent,
			childB: parent,
			childC: parent,
			parent: mainBranch,
		})
		var connector forgedomain.ProposalFinder = &testConnector{}
		args := proposallineage.ProposalStackLineageArgs{
			Connector:                Some(connector),
			CurrentBranch:            parent,
			Lineage:                  lineage,
			MainAndPerennialBranches: gitdomain.LocalBranchNames{mainBranch},
			Order:                    configdomain.OrderAsc,
		}

		tree, err := proposallineage.NewTree(args)

		must.NoError(t, err)
		must.NotNil(t, tree)
		must.True(t, tree.BranchToProposal[parent].IsSome())
		must.True(t, tree.BranchToProposal[childA].IsSome())
		must.True(t, tree.BranchToProposal[childB].IsSome())
		must.True(t, tree.BranchToProposal[childC].IsSome())
	})

	t.Run("handles multiple independent stacks", func(t *testing.T) {
		t.Parallel()
		mainBranch := gitdomain.NewLocalBranchName("main")
		stackA1 := gitdomain.NewLocalBranchName("stack-a-1")
		stackA2 := gitdomain.NewLocalBranchName("stack-a-2")
		stackB1 := gitdomain.NewLocalBranchName("stack-b-1")
		stackB2 := gitdomain.NewLocalBranchName("stack-b-2")
		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			stackA1: mainBranch,
			stackA2: stackA1,
			stackB1: mainBranch,
			stackB2: stackB1,
		})
		var connector forgedomain.ProposalFinder = &testConnector{}
		args := proposallineage.ProposalStackLineageArgs{
			Connector:                Some(connector),
			CurrentBranch:            stackA1,
			Lineage:                  lineage,
			MainAndPerennialBranches: gitdomain.LocalBranchNames{mainBranch},
		}

		tree, err := proposallineage.NewTree(args)

		must.NoError(t, err)
		must.NotNil(t, tree)
		// Should only include stack A in the tree
		must.True(t, tree.BranchToProposal[stackA1].IsSome())
		must.True(t, tree.BranchToProposal[stackA2].IsSome())
		// Stack B should not be included because it's not in current branch's stack
		must.True(t, tree.BranchToProposal[stackB1].IsNone())
		must.True(t, tree.BranchToProposal[stackB2].IsNone())
	})
}

func TestTreeWithMixedProposalAvailability(t *testing.T) {
	t.Parallel()

	t.Run("handles mix of branches with and without proposals", func(t *testing.T) {
		t.Parallel()
		mainBranch := gitdomain.NewLocalBranchName("main")
		withProposal := gitdomain.NewLocalBranchName("with-proposal")
		withoutProposal := gitdomain.NewLocalBranchName("no-proposal_branch")
		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			withProposal:    mainBranch,
			withoutProposal: withProposal,
		})
		var connector forgedomain.ProposalFinder = &testConnector{}
		args := proposallineage.ProposalStackLineageArgs{
			Connector:                Some(connector),
			CurrentBranch:            withProposal,
			Lineage:                  lineage,
			MainAndPerennialBranches: gitdomain.LocalBranchNames{mainBranch},
		}

		tree, err := proposallineage.NewTree(args)

		must.NoError(t, err)
		must.NotNil(t, tree)
		must.True(t, tree.BranchToProposal[withProposal].IsSome())
		must.True(t, tree.BranchToProposal[withoutProposal].IsNone())
	})
}
