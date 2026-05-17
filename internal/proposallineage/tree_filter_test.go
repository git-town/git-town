package proposallineage_test

import (
	"testing"

	"github.com/git-town/git-town/v23/internal/config/configdomain"
	"github.com/git-town/git-town/v23/internal/git/gitdomain"
	"github.com/git-town/git-town/v23/internal/proposallineage"
	. "github.com/git-town/git-town/v23/pkg/prelude"
	"github.com/git-town/git-town/v23/pkg/set"
	"github.com/shoenig/test/must"
)

func TestFilterTree(t *testing.T) {
	t.Parallel()

	t.Run("Exclude leaf node", func(t *testing.T) {
		t.Parallel()

		tree := proposallineage.TreeNode{
			Branch:        "main",
			LineageParent: None[gitdomain.LocalBranchName](),
			Children: proposallineage.Forest{
				proposallineage.TreeNode{
					Branch:        "feature",
					LineageParent: Some(gitdomain.LocalBranchName("main")),
					Children: proposallineage.Forest{
						proposallineage.TreeNode{
							Branch:        "prototype",
							LineageParent: Some(gitdomain.LocalBranchName("feature")),
							Children:      proposallineage.Forest{},
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
		excluded := set.New(configdomain.BranchTypePrototypeBranch)

		have := proposallineage.FilterTree(tree, branchTypes, excluded)

		want := proposallineage.Forest{
			proposallineage.TreeNode{
				Branch:        "main",
				LineageParent: None[gitdomain.LocalBranchName](),
				Children: proposallineage.Forest{
					proposallineage.TreeNode{
						Branch:        "feature",
						LineageParent: Some(gitdomain.LocalBranchName("main")),
						Children:      proposallineage.Forest{},
					},
				},
			},
		}

		must.Eq(t, want, have)
	})

	t.Run("Exclude middle node", func(t *testing.T) {
		t.Parallel()

		tree := proposallineage.TreeNode{
			Branch:        "main",
			LineageParent: None[gitdomain.LocalBranchName](),
			Children: proposallineage.Forest{
				proposallineage.TreeNode{
					Branch:        "feature-a",
					LineageParent: Some(gitdomain.LocalBranchName("main")),
					Children: proposallineage.Forest{
						proposallineage.TreeNode{
							Branch:        "prototype",
							LineageParent: Some(gitdomain.LocalBranchName("feature-a")),
							Children: proposallineage.Forest{
								proposallineage.TreeNode{
									Branch:        "feature-b",
									LineageParent: Some(gitdomain.LocalBranchName("prototype")),
									Children:      proposallineage.Forest{},
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
		excluded := set.New(configdomain.BranchTypePrototypeBranch)

		have := proposallineage.FilterTree(tree, branchTypes, excluded)

		want := proposallineage.Forest{
			proposallineage.TreeNode{
				Branch:        "main",
				LineageParent: None[gitdomain.LocalBranchName](),
				Children: proposallineage.Forest{
					proposallineage.TreeNode{
						Branch:        "feature-a",
						LineageParent: Some(gitdomain.LocalBranchName("main")),
						Children: proposallineage.Forest{
							proposallineage.TreeNode{
								Branch:        "feature-b",
								LineageParent: Some(gitdomain.LocalBranchName("prototype")),
								Children:      proposallineage.Forest{},
							},
						},
					},
				},
			},
		}

		must.Eq(t, want, have)
	})

	t.Run("Exclude middle node with multiple children", func(t *testing.T) {
		t.Parallel()

		tree := proposallineage.TreeNode{
			Branch:        "main",
			LineageParent: None[gitdomain.LocalBranchName](),
			Children: proposallineage.Forest{
				proposallineage.TreeNode{
					Branch:        "prototype",
					LineageParent: Some(gitdomain.LocalBranchName("main")),
					Children: proposallineage.Forest{
						proposallineage.TreeNode{
							Branch:        "feature-a",
							LineageParent: Some(gitdomain.LocalBranchName("prototype")),
							Children:      proposallineage.Forest{},
						},
						proposallineage.TreeNode{
							Branch:        "feature-b",
							LineageParent: Some(gitdomain.LocalBranchName("prototype")),
							Children:      proposallineage.Forest{},
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
		excluded := set.New(configdomain.BranchTypePrototypeBranch)

		have := proposallineage.FilterTree(tree, branchTypes, excluded)
		// This assertion explains why we used the "FOREST" terminology.
		// Filtering can POTENTIALLY lead to multiple stacks / TreeNodes
		// leading to main branch.
		want := proposallineage.Forest{
			proposallineage.TreeNode{
				Branch:        "main",
				LineageParent: None[gitdomain.LocalBranchName](),
				Children: proposallineage.Forest{
					proposallineage.TreeNode{
						Branch:        "feature-a",
						LineageParent: Some(gitdomain.LocalBranchName("prototype")),
						Children:      proposallineage.Forest{},
					},
					proposallineage.TreeNode{
						Branch:        "feature-b",
						LineageParent: Some(gitdomain.LocalBranchName("prototype")),
						Children:      proposallineage.Forest{},
					},
				},
			},
		}

		must.Eq(t, want, have)
	})

	t.Run("Exclude root node", func(t *testing.T) {
		t.Parallel()

		tree := proposallineage.TreeNode{
			Branch:        "main",
			LineageParent: None[gitdomain.LocalBranchName](),
			Children: proposallineage.Forest{
				proposallineage.TreeNode{
					Branch:        "feature-a",
					LineageParent: Some(gitdomain.LocalBranchName("main")),
					Children: proposallineage.Forest{
						proposallineage.TreeNode{
							Branch:        "feature-b",
							LineageParent: Some(gitdomain.LocalBranchName("feature-a")),
							Children:      proposallineage.Forest{},
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
		excluded := set.New(configdomain.BranchTypeMainBranch)

		have := proposallineage.FilterTree(tree, branchTypes, excluded)

		want := proposallineage.Forest{
			proposallineage.TreeNode{
				Branch:        "feature-a",
				LineageParent: Some(gitdomain.LocalBranchName("main")),
				Children: proposallineage.Forest{
					proposallineage.TreeNode{
						Branch:        "feature-b",
						LineageParent: Some(gitdomain.LocalBranchName("feature-a")),
						Children:      proposallineage.Forest{},
					},
				},
			},
		}

		must.Eq(t, want, have)
	})
}
