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
type testConnector struct{}

func (self *testConnector) FindProposal(source, target gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
	if strings.Contains(source.String(), "no_proposal") {
		return None[forgedomain.Proposal](), nil
	}
	return Some(forgedomain.Proposal{
		Data: forgedomain.ProposalData{
			Title: gitdomain.ProposalTitle(fmt.Sprintf("proposal from %s to %s", source, target)),
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

	t.Run("branch in the middle of a stack", func(t *testing.T) {
		t.Parallel()
		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			"feature-a": "main",
			"feature-b": "feature-a",
			"feature-c": "feature-a",
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
	})

	t.Run("when on a leaf branch, it prints the leaf branch first and then the other siblings", func(t *testing.T) {
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

	t.Run("branch in the middle of a long lineage", func(t *testing.T) {
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

	t.Run("no connector", func(t *testing.T) {
		t.Parallel()
		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			"feature-a": "main",
			"feature-b": "feature-a",
			"feature-c": "feature-a",
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
				"feature-c": None[forgedomain.Proposal](),
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
							{
								Branch:     "feature-c",
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

	t.Run("creates tree without connector", func(t *testing.T) {
		t.Parallel()
		mainBranch := gitdomain.NewLocalBranchName("main")
		featureBranch := gitdomain.NewLocalBranchName("feature")
		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			featureBranch: mainBranch,
		})
		args := proposallineage.ProposalStackLineageArgs{
			Connector:                None[forgedomain.ProposalFinder](),
			CurrentBranch:            featureBranch,
			Lineage:                  lineage,
			MainAndPerennialBranches: gitdomain.LocalBranchNames{mainBranch},
		}

		tree, err := proposallineage.NewTree(args)

		must.NoError(t, err)
		must.NotNil(t, tree)
		must.True(t, tree.BranchToProposal[featureBranch].IsNone())
	})

	t.Run("handles branches with no proposals", func(t *testing.T) {
		t.Parallel()
		mainBranch := gitdomain.NewLocalBranchName("main")
		noproposalBranch := gitdomain.NewLocalBranchName("no_proposal_branch")
		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			noproposalBranch: mainBranch,
		})
		var connector forgedomain.ProposalFinder = &testConnector{}
		args := proposallineage.ProposalStackLineageArgs{
			Connector:                Some(connector),
			CurrentBranch:            noproposalBranch,
			Lineage:                  lineage,
			MainAndPerennialBranches: gitdomain.LocalBranchNames{mainBranch},
		}

		tree, err := proposallineage.NewTree(args)

		must.NoError(t, err)
		must.NotNil(t, tree)
		must.True(t, tree.BranchToProposal[noproposalBranch].IsNone())
	})

	t.Run("handles error from proposal finder", func(t *testing.T) {
		t.Parallel()
		mainBranch := gitdomain.NewLocalBranchName("main")
		featureBranch := gitdomain.NewLocalBranchName("feature")
		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			featureBranch: mainBranch,
		})
		var connector forgedomain.ProposalFinder = &failingConnector{}
		args := proposallineage.ProposalStackLineageArgs{
			Connector:                Some(connector),
			CurrentBranch:            featureBranch,
			Lineage:                  lineage,
			MainAndPerennialBranches: gitdomain.LocalBranchNames{mainBranch},
		}

		tree, err := proposallineage.NewTree(args)

		must.Error(t, err)
		must.NotNil(t, tree)
	})
}

func TestTreeCachingBehavior(t *testing.T) {
	t.Parallel()

	t.Run("caches proposals across rebuild", func(t *testing.T) {
		t.Parallel()
		mainBranch := gitdomain.NewLocalBranchName("main")
		featureA := gitdomain.NewLocalBranchName("feature-a")
		featureB := gitdomain.NewLocalBranchName("feature-b")
		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			featureA: mainBranch,
			featureB: featureA,
		})
		var connector forgedomain.ProposalFinder = &testConnector{}
		args := proposallineage.ProposalStackLineageArgs{
			Connector:                Some(connector),
			CurrentBranch:            featureA,
			Lineage:                  lineage,
			MainAndPerennialBranches: gitdomain.LocalBranchNames{mainBranch},
		}

		tree, err := proposallineage.NewTree(args)
		must.NoError(t, err)

		// Get proposals from first build
		proposalA1 := tree.BranchToProposal[featureA]
		proposalB1 := tree.BranchToProposal[featureB]

		// Rebuild
		err = tree.Rebuild(args)
		must.NoError(t, err)

		// Proposals should be the same (reused from cache)
		proposalA2 := tree.BranchToProposal[featureA]
		proposalB2 := tree.BranchToProposal[featureB]

		must.Eq(t, proposalA1, proposalA2)
		must.Eq(t, proposalB1, proposalB2)
	})

	t.Run("fetches new proposals for new branches after rebuild", func(t *testing.T) {
		t.Parallel()
		mainBranch := gitdomain.NewLocalBranchName("main")
		featureA := gitdomain.NewLocalBranchName("feature-a")
		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			featureA: mainBranch,
		})
		var connector forgedomain.ProposalFinder = &testConnector{}
		args := proposallineage.ProposalStackLineageArgs{
			Connector:                Some(connector),
			CurrentBranch:            featureA,
			Lineage:                  lineage,
			MainAndPerennialBranches: gitdomain.LocalBranchNames{mainBranch},
		}

		tree, err := proposallineage.NewTree(args)
		must.NoError(t, err)

		// featureB is not in cache
		featureB := gitdomain.NewLocalBranchName("feature-b")
		must.True(t, tree.BranchToProposal[featureB].IsNone())

		// Add featureB to lineage and rebuild
		newLineage := configdomain.NewLineageWith(configdomain.LineageData{
			featureA: mainBranch,
			featureB: featureA,
		})
		newArgs := proposallineage.ProposalStackLineageArgs{
			Connector:                Some(connector),
			CurrentBranch:            featureA,
			Lineage:                  newLineage,
			MainAndPerennialBranches: gitdomain.LocalBranchNames{mainBranch},
		}

		err = tree.Rebuild(newArgs)
		must.NoError(t, err)

		// featureB should now have a proposal
		must.True(t, tree.BranchToProposal[featureB].IsSome())
	})
}

func TestTreeRebuild(t *testing.T) {
	t.Parallel()

	t.Run("handles error during rebuild", func(t *testing.T) {
		t.Parallel()
		mainBranch := gitdomain.NewLocalBranchName("main")
		featureA := gitdomain.NewLocalBranchName("feature-a")
		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			featureA: mainBranch,
		})
		var connector forgedomain.ProposalFinder = &testConnector{}
		args := proposallineage.ProposalStackLineageArgs{
			Connector:                Some(connector),
			CurrentBranch:            featureA,
			Lineage:                  lineage,
			MainAndPerennialBranches: gitdomain.LocalBranchNames{mainBranch},
		}

		tree, err := proposallineage.NewTree(args)
		must.NoError(t, err)

		// Rebuild with error-inducing connector
		var errorConnector forgedomain.ProposalFinder = &failingConnector{}
		featureB := gitdomain.NewLocalBranchName("feature-b")
		newLineage := configdomain.NewLineageWith(configdomain.LineageData{
			featureA: mainBranch,
			featureB: featureA,
		})
		errorArgs := proposallineage.ProposalStackLineageArgs{
			Connector:                Some(errorConnector),
			CurrentBranch:            featureB,
			Lineage:                  newLineage,
			MainAndPerennialBranches: gitdomain.LocalBranchNames{mainBranch},
		}

		err = tree.Rebuild(errorArgs)
		must.Error(t, err)
	})

	t.Run("rebuilds tree and reuses cached proposals", func(t *testing.T) {
		t.Parallel()
		mainBranch := gitdomain.NewLocalBranchName("main")
		featureA := gitdomain.NewLocalBranchName("feature-a")
		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			featureA: mainBranch,
		})
		var connector forgedomain.ProposalFinder = &testConnector{}
		args := proposallineage.ProposalStackLineageArgs{
			Connector:                Some(connector),
			CurrentBranch:            featureA,
			Lineage:                  lineage,
			MainAndPerennialBranches: gitdomain.LocalBranchNames{mainBranch},
		}

		tree, err := proposallineage.NewTree(args)
		must.NoError(t, err)
		originalProposal := tree.BranchToProposal[featureA]

		// Rebuild with same args
		err = tree.Rebuild(args)
		must.NoError(t, err)

		// Should reuse the cached proposal
		must.Eq(t, originalProposal, tree.BranchToProposal[featureA])
	})

	t.Run("rebuilds tree with updated lineage", func(t *testing.T) {
		t.Parallel()
		mainBranch := gitdomain.NewLocalBranchName("main")
		featureA := gitdomain.NewLocalBranchName("feature-a")
		featureB := gitdomain.NewLocalBranchName("feature-b")
		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			featureA: mainBranch,
			featureB: featureA,
		})
		var connector forgedomain.ProposalFinder = &testConnector{}
		args := proposallineage.ProposalStackLineageArgs{
			Connector:                Some(connector),
			CurrentBranch:            featureA,
			Lineage:                  lineage,
			MainAndPerennialBranches: gitdomain.LocalBranchNames{mainBranch},
		}

		tree, err := proposallineage.NewTree(args)
		must.NoError(t, err)

		// Now update the lineage and rebuild
		featureC := gitdomain.NewLocalBranchName("feature-c")
		newLineage := configdomain.NewLineageWith(configdomain.LineageData{
			featureA: mainBranch,
			featureB: featureA,
			featureC: featureB,
		})
		newArgs := proposallineage.ProposalStackLineageArgs{
			Connector:                Some(connector),
			CurrentBranch:            featureB,
			Lineage:                  newLineage,
			MainAndPerennialBranches: gitdomain.LocalBranchNames{mainBranch},
		}

		err = tree.Rebuild(newArgs)

		must.NoError(t, err)
		must.NotNil(t, tree.Node)
		// BranchToProposal should retain cached proposals from first build
		must.True(t, tree.BranchToProposal[featureA].IsSome())
		must.True(t, tree.BranchToProposal[featureB].IsSome())
		must.True(t, tree.BranchToProposal[featureC].IsSome())
	})
}

func TestTreeWithComplexLineages(t *testing.T) {
	t.Parallel()

	t.Run("handles current branch on perennial", func(t *testing.T) {
		t.Parallel()
		mainBranch := gitdomain.NewLocalBranchName("main")
		featureBranch := gitdomain.NewLocalBranchName("feature")
		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			featureBranch: mainBranch,
		})
		var connector forgedomain.ProposalFinder = &testConnector{}
		args := proposallineage.ProposalStackLineageArgs{
			Connector:                Some(connector),
			CurrentBranch:            mainBranch, // Current is perennial
			Lineage:                  lineage,
			MainAndPerennialBranches: gitdomain.LocalBranchNames{mainBranch},
		}

		tree, err := proposallineage.NewTree(args)

		must.NoError(t, err)
		must.NotNil(t, tree)
		// When current branch is perennial, feature branches are not processed
		// because the tree only builds from ancestors and descendants of current branch
		must.True(t, tree.BranchToProposal[featureBranch].IsNone())
	})

	t.Run("handles deep nesting", func(t *testing.T) {
		t.Parallel()
		mainBranch := gitdomain.NewLocalBranchName("main")
		branches := make([]gitdomain.LocalBranchName, 10)
		lineageData := configdomain.LineageData{}

		branches[0] = gitdomain.NewLocalBranchName("level-0")
		lineageData[branches[0]] = mainBranch

		for i := 1; i < 10; i++ {
			branches[i] = gitdomain.NewLocalBranchName(fmt.Sprintf("level-%d", i))
			lineageData[branches[i]] = branches[i-1]
		}

		lineage := configdomain.NewLineageWith(lineageData)
		var connector forgedomain.ProposalFinder = &testConnector{}
		args := proposallineage.ProposalStackLineageArgs{
			Connector:                Some(connector),
			CurrentBranch:            branches[4], // Middle of the chain
			Lineage:                  lineage,
			MainAndPerennialBranches: gitdomain.LocalBranchNames{mainBranch},
		}

		tree, err := proposallineage.NewTree(args)

		must.NoError(t, err)
		must.NotNil(t, tree)
		// All branches in the chain should have proposals
		for i := range 10 {
			must.True(t, tree.BranchToProposal[branches[i]].IsSome())
		}
	})

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
		withoutProposal := gitdomain.NewLocalBranchName("no_proposal_branch")
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
