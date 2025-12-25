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

func TestNewTree(t *testing.T) {
	t.Parallel()

	t.Run("connector returns errors", func(t *testing.T) {
		t.Parallel()
		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			"feature-a": "main",
		})
		var connector forgedomain.ProposalFinder = &failingFinder{}
		_, err := proposallineage.NewTree(proposallineage.ProposalStackLineageArgs{
			Connector:                Some(connector),
			CurrentBranch:            "feature-a",
			Lineage:                  lineage,
			MainAndPerennialBranches: gitdomain.LocalBranchNames{"main"},
		})
		must.Error(t, err)
	})

	t.Run("feature branch", func(t *testing.T) {
		t.Parallel()
		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			"feature-a": "main",
		})
		var connector forgedomain.ProposalFinder = &testFinder{}
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

	t.Run("feature branch in a long lineage", func(t *testing.T) {
		t.Parallel()
		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			"feature-a": "main",
			"feature-b": "feature-a",
			"feature-c": "feature-b",
			"feature-d": "feature-c",
		})
		var connector forgedomain.ProposalFinder = &testFinder{}
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

	t.Run("feature branch with multiple children", func(t *testing.T) {
		t.Parallel()
		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			"feature-a": "main",
			"feature-b": "feature-a",
			"feature-c": "feature-a",
		})
		connector := testFinder{}
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

	t.Run("leaf branch with siblings", func(t *testing.T) {
		t.Parallel()
		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			"feature-a": "main",
			"feature-b": "feature-a",
			"feature-c": "feature-a",
			"feature-d": "feature-a",
		})
		var connector forgedomain.ProposalFinder = &testFinder{}
		have, err := proposallineage.NewTree(proposallineage.ProposalStackLineageArgs{
			Connector:                Some(connector),
			CurrentBranch:            "feature-d",
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
				"feature-d": Some(forgedomain.Proposal{
					Data: forgedomain.ProposalData{
						Title: "proposal from feature-d to feature-a",
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
								Branch:     "feature-d",
								ChildNodes: []*proposallineage.TreeNode{},
								Proposal: Some(forgedomain.Proposal{
									Data: forgedomain.ProposalData{
										Title: "proposal from feature-d to feature-a",
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

	t.Run("perennial branch at the root", func(t *testing.T) {
		t.Parallel()
		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			"feature-a": "main",
			"feature-b": "feature-a",
		})
		var connector forgedomain.ProposalFinder = &testFinder{}
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

	t.Run("several independent stacks", func(t *testing.T) {
		t.Parallel()
		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			"feature-A1": "main",
			"feature-A2": "feature-A1",
			"feature-B1": "main",
			"feature-B2": "feature-B1",
		})
		var connector forgedomain.ProposalFinder = &testFinder{}
		have, err := proposallineage.NewTree(proposallineage.ProposalStackLineageArgs{
			Connector:                Some(connector),
			CurrentBranch:            "feature-A1",
			Lineage:                  lineage,
			MainAndPerennialBranches: gitdomain.LocalBranchNames{"main"},
		})
		want := &proposallineage.Tree{
			BranchToProposal: map[gitdomain.LocalBranchName]Option[forgedomain.Proposal]{
				"feature-A1": Some(forgedomain.Proposal{
					Data: forgedomain.ProposalData{
						Title: "proposal from feature-A1 to main",
					},
				}),
				"feature-A2": Some(forgedomain.Proposal{
					Data: forgedomain.ProposalData{
						Title: "proposal from feature-A2 to feature-A1",
					},
				}),
			},
			Node: &proposallineage.TreeNode{
				Branch: "main",
				ChildNodes: []*proposallineage.TreeNode{
					{
						Branch: "feature-A1",
						ChildNodes: []*proposallineage.TreeNode{
							{
								Branch:     "feature-A2",
								ChildNodes: []*proposallineage.TreeNode{},
								Proposal: Some(forgedomain.Proposal{
									Data: forgedomain.ProposalData{
										Title: "proposal from feature-A2 to feature-A1",
									},
								}),
							},
						},
						Proposal: Some(forgedomain.Proposal{
							Data: forgedomain.ProposalData{
								Title: "proposal from feature-A1 to main",
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
		var connector forgedomain.ProposalFinder = &testFinder{}
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
		connector := testFinder{}
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
		connector = testFinder{}
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
		connector := testFinder{}
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
		var errorConnector forgedomain.ProposalFinder = &failingFinder{}
		err = tree.Rebuild(proposallineage.ProposalStackLineageArgs{
			Connector:                Some(errorConnector),
			CurrentBranch:            "feature-b",
			Lineage:                  lineage2,
			MainAndPerennialBranches: gitdomain.LocalBranchNames{"main"},
		})
		must.Error(t, err)
	})
}
