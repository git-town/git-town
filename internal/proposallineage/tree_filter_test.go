package proposallineage_test

import (
	"testing"

	"github.com/git-town/git-town/v23/internal/config/configdomain"
	"github.com/git-town/git-town/v23/internal/git/gitdomain"
	"github.com/git-town/git-town/v23/internal/proposallineage"
	. "github.com/git-town/git-town/v23/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestFilterTree(t *testing.T) {
	t.Parallel()

	t.Run("Exclude leaf node", func(t *testing.T) {
		t.Parallel()

		tree := proposallineage.TreeNode{
			Branch:        "main",
			LineageParent: None[gitdomain.LocalBranchName](),
			Children: proposallineage.TreeNodes{
				proposallineage.TreeNode{
					Branch:        "feature",
					LineageParent: Some(gitdomain.LocalBranchName("main")),
					Children: proposallineage.TreeNodes{
						proposallineage.TreeNode{
							Branch:        "prototype",
							LineageParent: Some(gitdomain.LocalBranchName("feature")),
							Children:      proposallineage.TreeNodes{},
						},
					},
				},
			},
		}
		branchTypes := configdomain.BranchesAndTypes{
			"main":      configdomain.BranchTypeMainBranch,
			"feature":   configdomain.BranchTypeFeatureBranch,
			"prototype": configdomain.BranchTypePrototypeBranch,
		}
		excluded := configdomain.NewProposalBreadcrumbExcludeBranches(configdomain.BranchTypePrototypeBranch)

		have := proposallineage.FilterTree(treeWithBranchTypes(tree, branchTypes), excluded)

		want := proposallineage.TreeNodes{
			proposallineage.TreeNode{
				Branch:        "main",
				LineageParent: None[gitdomain.LocalBranchName](),
				Children: proposallineage.TreeNodes{
					proposallineage.TreeNode{
						Branch:        "feature",
						LineageParent: Some(gitdomain.LocalBranchName("main")),
						Children:      proposallineage.TreeNodes{},
					},
				},
			},
		}

		must.Eq(t, treeNodesWithBranchTypes(want, branchTypes), have)
	})

	t.Run("Exclude middle node", func(t *testing.T) {
		t.Parallel()

		tree := proposallineage.TreeNode{
			Branch:        "main",
			LineageParent: None[gitdomain.LocalBranchName](),
			Children: proposallineage.TreeNodes{
				proposallineage.TreeNode{
					Branch:        "feature-a",
					LineageParent: Some(gitdomain.LocalBranchName("main")),
					Children: proposallineage.TreeNodes{
						proposallineage.TreeNode{
							Branch:        "prototype",
							LineageParent: Some(gitdomain.LocalBranchName("feature-a")),
							Children: proposallineage.TreeNodes{
								proposallineage.TreeNode{
									Branch:        "feature-b",
									LineageParent: Some(gitdomain.LocalBranchName("prototype")),
									Children:      proposallineage.TreeNodes{},
								},
							},
						},
					},
				},
			},
		}
		branchTypes := configdomain.BranchesAndTypes{
			"main":      configdomain.BranchTypeMainBranch,
			"feature-a": configdomain.BranchTypeFeatureBranch,
			"feature-b": configdomain.BranchTypeFeatureBranch,
			"prototype": configdomain.BranchTypePrototypeBranch,
		}
		excluded := configdomain.NewProposalBreadcrumbExcludeBranches(configdomain.BranchTypePrototypeBranch)

		have := proposallineage.FilterTree(treeWithBranchTypes(tree, branchTypes), excluded)

		want := proposallineage.TreeNodes{
			proposallineage.TreeNode{
				Branch:        "main",
				LineageParent: None[gitdomain.LocalBranchName](),
				Children: proposallineage.TreeNodes{
					proposallineage.TreeNode{
						Branch:        "feature-a",
						LineageParent: Some(gitdomain.LocalBranchName("main")),
						Children: proposallineage.TreeNodes{
							proposallineage.TreeNode{
								Branch:        "feature-b",
								LineageParent: Some(gitdomain.LocalBranchName("prototype")),
								Children:      proposallineage.TreeNodes{},
							},
						},
					},
				},
			},
		}

		must.Eq(t, treeNodesWithBranchTypes(want, branchTypes), have)
	})

	t.Run("Exclude middle node with multiple children", func(t *testing.T) {
		t.Parallel()

		tree := proposallineage.TreeNode{
			Branch:        "main",
			LineageParent: None[gitdomain.LocalBranchName](),
			Children: proposallineage.TreeNodes{
				proposallineage.TreeNode{
					Branch:        "prototype",
					LineageParent: Some(gitdomain.LocalBranchName("main")),
					Children: proposallineage.TreeNodes{
						proposallineage.TreeNode{
							Branch:        "feature-a",
							LineageParent: Some(gitdomain.LocalBranchName("prototype")),
							Children:      proposallineage.TreeNodes{},
						},
						proposallineage.TreeNode{
							Branch:        "feature-b",
							LineageParent: Some(gitdomain.LocalBranchName("prototype")),
							Children:      proposallineage.TreeNodes{},
						},
					},
				},
			},
		}
		branchTypes := configdomain.BranchesAndTypes{
			"main":      configdomain.BranchTypeMainBranch,
			"feature-a": configdomain.BranchTypeFeatureBranch,
			"feature-b": configdomain.BranchTypeFeatureBranch,
			"prototype": configdomain.BranchTypePrototypeBranch,
		}
		excluded := configdomain.NewProposalBreadcrumbExcludeBranches(configdomain.BranchTypePrototypeBranch)

		have := proposallineage.FilterTree(treeWithBranchTypes(tree, branchTypes), excluded)
		// This assertion explains why we used the "FOREST" terminology.
		// Filtering can POTENTIALLY lead to multiple stacks / TreeNodes
		// leading to main branch.
		want := proposallineage.TreeNodes{
			proposallineage.TreeNode{
				Branch:        "main",
				LineageParent: None[gitdomain.LocalBranchName](),
				Children: proposallineage.TreeNodes{
					proposallineage.TreeNode{
						Branch:        "feature-a",
						LineageParent: Some(gitdomain.LocalBranchName("prototype")),
						Children:      proposallineage.TreeNodes{},
					},
					proposallineage.TreeNode{
						Branch:        "feature-b",
						LineageParent: Some(gitdomain.LocalBranchName("prototype")),
						Children:      proposallineage.TreeNodes{},
					},
				},
			},
		}

		must.Eq(t, treeNodesWithBranchTypes(want, branchTypes), have)
	})

	t.Run("Exclude root node", func(t *testing.T) {
		t.Parallel()

		tree := proposallineage.TreeNode{
			Branch:        "main",
			LineageParent: None[gitdomain.LocalBranchName](),
			Children: proposallineage.TreeNodes{
				proposallineage.TreeNode{
					Branch:        "feature-a",
					LineageParent: Some(gitdomain.LocalBranchName("main")),
					Children: proposallineage.TreeNodes{
						proposallineage.TreeNode{
							Branch:        "feature-b",
							LineageParent: Some(gitdomain.LocalBranchName("feature-a")),
							Children:      proposallineage.TreeNodes{},
						},
					},
				},
			},
		}
		branchTypes := configdomain.BranchesAndTypes{
			"main":      configdomain.BranchTypeMainBranch,
			"feature-a": configdomain.BranchTypeFeatureBranch,
			"feature-b": configdomain.BranchTypeFeatureBranch,
		}
		excluded := configdomain.NewProposalBreadcrumbExcludeBranches(configdomain.BranchTypeMainBranch)

		have := proposallineage.FilterTree(treeWithBranchTypes(tree, branchTypes), excluded)

		want := proposallineage.TreeNodes{
			proposallineage.TreeNode{
				Branch:        "feature-a",
				LineageParent: Some(gitdomain.LocalBranchName("main")),
				Children: proposallineage.TreeNodes{
					proposallineage.TreeNode{
						Branch:        "feature-b",
						LineageParent: Some(gitdomain.LocalBranchName("feature-a")),
						Children:      proposallineage.TreeNodes{},
					},
				},
			},
		}

		must.Eq(t, treeNodesWithBranchTypes(want, branchTypes), have)
	})
}

func treeNodesWithBranchTypes(nodes proposallineage.TreeNodes, branchTypes configdomain.BranchesAndTypes) proposallineage.TreeNodes {
	for index, node := range nodes {
		nodes[index] = treeWithBranchTypes(node, branchTypes)
	}
	return nodes
}

func treeWithBranchTypes(tree proposallineage.TreeNode, branchTypes configdomain.BranchesAndTypes) proposallineage.TreeNode {
	tree.BranchType = branchTypes[tree.Branch]
	for index, child := range tree.Children {
		tree.Children[index] = treeWithBranchTypes(child, branchTypes)
	}
	return tree
}
