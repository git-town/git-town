package proposallineage_test

import (
	"testing"

	"github.com/git-town/git-town/v23/internal/config/configdomain"
	"github.com/git-town/git-town/v23/internal/git/gitdomain"
	"github.com/git-town/git-town/v23/internal/proposallineage"
	. "github.com/git-town/git-town/v23/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestTreeNode(t *testing.T) {
	t.Parallel()

	t.Run("BranchCount", func(t *testing.T) {
		t.Parallel()
		t.Run("deep tree", func(t *testing.T) {
			t.Parallel()
			tree := proposallineage.TreeNode{
				Branch:        "main",
				LineageParent: None[gitdomain.LocalBranchName](),
				Children: []proposallineage.TreeNode{
					{Branch: "branch-1", LineageParent: Some(gitdomain.LocalBranchName("main"))},
					{Branch: "branch-2", LineageParent: Some(gitdomain.LocalBranchName("main")), Children: []proposallineage.TreeNode{
						{Branch: "branch-2a", LineageParent: Some(gitdomain.LocalBranchName("branch-2")), Children: []proposallineage.TreeNode{
							{Branch: "branch-2a1", LineageParent: Some(gitdomain.LocalBranchName("branch-2a"))},
						}},
					}},
				},
			}
			have := tree.BranchCount()
			must.EqOp(t, 5, have)
		})
		t.Run("single branch", func(t *testing.T) {
			t.Parallel()
			tree := proposallineage.TreeNode{
				Branch:        "main",
				LineageParent: None[gitdomain.LocalBranchName](),
				Children: []proposallineage.TreeNode{
					{Branch: "branch-1", LineageParent: Some(gitdomain.LocalBranchName("main"))},
				},
			}
			have := tree.BranchCount()
			must.EqOp(t, 2, have)
		})
	})

	t.Run("CalculateTree", func(t *testing.T) {
		t.Parallel()
		t.Run("branch in a long lineage", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.NewLineageWith(configdomain.LineageData{
				"feature-a": "main",
				"feature-b": "feature-a",
				"feature-c": "feature-b",
				"feature-d": "feature-c",
				"feature-e": "feature-d",
			})
			have := proposallineage.CalculateTree("feature-b", lineage, configdomain.OrderAsc, configdomain.BranchesAndTypes{})
			want := proposallineage.TreeNode{
				Branch:        "main",
				LineageParent: None[gitdomain.LocalBranchName](),
				Children: []proposallineage.TreeNode{
					{
						Branch:        "feature-a",
						LineageParent: Some(gitdomain.LocalBranchName("main")),
						Children: []proposallineage.TreeNode{
							{
								Branch:        "feature-b",
								LineageParent: Some(gitdomain.LocalBranchName("feature-a")),
								Children: []proposallineage.TreeNode{
									{
										Branch:        "feature-c",
										LineageParent: Some(gitdomain.LocalBranchName("feature-b")),
										Children: []proposallineage.TreeNode{
											{
												Branch:        "feature-d",
												LineageParent: Some(gitdomain.LocalBranchName("feature-c")),
												Children: []proposallineage.TreeNode{
													{
														Branch:        "feature-e",
														LineageParent: Some(gitdomain.LocalBranchName("feature-d")),
														Children:      []proposallineage.TreeNode{},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			}
			must.Eq(t, want, have)
		})

		t.Run("branch with multiple descendent lineages", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.NewLineageWith(configdomain.LineageData{
				"feature-a":   "main",
				"feature-b1":  "feature-a",
				"feature-b1a": "feature-b1",
				"feature-b1b": "feature-b1",
				"feature-b2":  "feature-a",
				"feature-b2a": "feature-b2",
				"feature-b2b": "feature-b2",
			})
			have := proposallineage.CalculateTree("feature-a", lineage, configdomain.OrderAsc, configdomain.BranchesAndTypes{})
			want := proposallineage.TreeNode{
				Branch:        "main",
				LineageParent: None[gitdomain.LocalBranchName](),
				Children: []proposallineage.TreeNode{
					{
						Branch:        "feature-a",
						LineageParent: Some(gitdomain.LocalBranchName("main")),
						Children: []proposallineage.TreeNode{
							{
								Branch:        "feature-b1",
								LineageParent: Some(gitdomain.LocalBranchName("feature-a")),
								Children: []proposallineage.TreeNode{
									{
										Branch:        "feature-b1a",
										LineageParent: Some(gitdomain.LocalBranchName("feature-b1")),
										Children:      []proposallineage.TreeNode{},
									},
									{
										Branch:        "feature-b1b",
										LineageParent: Some(gitdomain.LocalBranchName("feature-b1")),
										Children:      []proposallineage.TreeNode{},
									},
								},
							},
							{
								Branch:        "feature-b2",
								LineageParent: Some(gitdomain.LocalBranchName("feature-a")),
								Children: []proposallineage.TreeNode{
									{
										Branch:        "feature-b2a",
										LineageParent: Some(gitdomain.LocalBranchName("feature-b2")),
										Children:      []proposallineage.TreeNode{},
									},
									{
										Branch:        "feature-b2b",
										LineageParent: Some(gitdomain.LocalBranchName("feature-b2")),
										Children:      []proposallineage.TreeNode{},
									},
								},
							},
						},
					},
				},
			}
			must.Eq(t, want, have)
		})

		t.Run("ignore independent lineages", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.NewLineageWith(configdomain.LineageData{
				"feature-a":  "main",
				"feature-a1": "feature-a",
				"feature-b":  "main",
				"feature-b1": "feature-b",
			})
			have := proposallineage.CalculateTree("feature-a", lineage, configdomain.OrderAsc, configdomain.BranchesAndTypes{})
			want := proposallineage.TreeNode{
				Branch:        "main",
				LineageParent: None[gitdomain.LocalBranchName](),
				Children: []proposallineage.TreeNode{
					{
						Branch:        "feature-a",
						LineageParent: Some(gitdomain.LocalBranchName("main")),
						Children: []proposallineage.TreeNode{
							{
								Branch:        "feature-a1",
								LineageParent: Some(gitdomain.LocalBranchName("feature-a")),
								Children:      []proposallineage.TreeNode{},
							},
						},
					},
				},
			}
			must.Eq(t, want, have)
		})

		t.Run("order descending", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.NewLineageWith(configdomain.LineageData{
				"feature-a":   "main",
				"feature-b1":  "feature-a",
				"feature-b1a": "feature-b1",
				"feature-b1b": "feature-b1",
				"feature-b2":  "feature-a",
				"feature-b2a": "feature-b2",
				"feature-b2b": "feature-b2",
			})
			have := proposallineage.CalculateTree("feature-a", lineage, configdomain.OrderDesc, configdomain.BranchesAndTypes{})
			want := proposallineage.TreeNode{
				Branch:        "main",
				LineageParent: None[gitdomain.LocalBranchName](),
				Children: []proposallineage.TreeNode{
					{
						Branch:        "feature-a",
						LineageParent: Some(gitdomain.LocalBranchName("main")),
						Children: []proposallineage.TreeNode{
							{
								Branch:        "feature-b2",
								LineageParent: Some(gitdomain.LocalBranchName("feature-a")),
								Children: []proposallineage.TreeNode{
									{
										Branch:        "feature-b2b",
										LineageParent: Some(gitdomain.LocalBranchName("feature-b2")),
										Children:      []proposallineage.TreeNode{},
									},
									{
										Branch:        "feature-b2a",
										LineageParent: Some(gitdomain.LocalBranchName("feature-b2")),
										Children:      []proposallineage.TreeNode{},
									},
								},
							},
							{
								Branch:        "feature-b1",
								LineageParent: Some(gitdomain.LocalBranchName("feature-a")),
								Children: []proposallineage.TreeNode{
									{
										Branch:        "feature-b1b",
										LineageParent: Some(gitdomain.LocalBranchName("feature-b1")),
										Children:      []proposallineage.TreeNode{},
									},
									{
										Branch:        "feature-b1a",
										LineageParent: Some(gitdomain.LocalBranchName("feature-b1")),
										Children:      []proposallineage.TreeNode{},
									},
								},
							},
						},
					},
				},
			}
			must.Eq(t, want, have)
		})

		t.Run("perennial branch", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.NewLineage()
			have := proposallineage.CalculateTree("main", lineage, configdomain.OrderAsc, configdomain.BranchesAndTypes{})
			want := proposallineage.TreeNode{
				Branch:        "main",
				LineageParent: None[gitdomain.LocalBranchName](),
				Children:      []proposallineage.TreeNode{},
			}
			must.Eq(t, want, have)
		})
	})
}
