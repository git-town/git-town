package proposallineage_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/proposallineage"
	"github.com/shoenig/test/must"
)

func TestCalculateTree(t *testing.T) {
	t.Parallel()

	t.Run("BranchCount", func(t *testing.T) {
		t.Parallel()
		tree := proposallineage.TreeNode{
			Branch: "main",
			Children: []proposallineage.TreeNode{
				{Branch: "branch-1"},
				{Branch: "branch-2", Children: []proposallineage.TreeNode{
					{Branch: "branch-2a", Children: []proposallineage.TreeNode{
						{Branch: "branch-2a1"},
					}},
				}},
			},
		}
		have := tree.BranchCount()
		must.EqOp(t, 5, have)
	})

	t.Run("CalculateTree", func(t *testing.T) {
		t.Run("branch in a long lineage", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.NewLineageWith(configdomain.LineageData{
				"feature-a": "main",
				"feature-b": "feature-a",
				"feature-c": "feature-b",
				"feature-d": "feature-c",
				"feature-e": "feature-d",
			})
			have := proposallineage.CalculateTree("feature-b", lineage, configdomain.OrderAsc)
			want := proposallineage.TreeNode{
				Branch: "main",
				Children: []proposallineage.TreeNode{
					{
						Branch: "feature-a",
						Children: []proposallineage.TreeNode{
							{
								Branch: "feature-b",
								Children: []proposallineage.TreeNode{
									{
										Branch: "feature-c",
										Children: []proposallineage.TreeNode{
											{
												Branch: "feature-d",
												Children: []proposallineage.TreeNode{
													{
														Branch:   "feature-e",
														Children: []proposallineage.TreeNode{},
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
			have := proposallineage.CalculateTree("feature-a", lineage, configdomain.OrderAsc)
			want := proposallineage.TreeNode{
				Branch: "main",
				Children: []proposallineage.TreeNode{
					{
						Branch: "feature-a",
						Children: []proposallineage.TreeNode{
							{
								Branch: "feature-b1",
								Children: []proposallineage.TreeNode{
									{
										Branch:   "feature-b1a",
										Children: []proposallineage.TreeNode{},
									},
									{
										Branch:   "feature-b1b",
										Children: []proposallineage.TreeNode{},
									},
								},
							},
							{
								Branch: "feature-b2",
								Children: []proposallineage.TreeNode{
									{
										Branch:   "feature-b2a",
										Children: []proposallineage.TreeNode{},
									},
									{
										Branch:   "feature-b2b",
										Children: []proposallineage.TreeNode{},
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
			have := proposallineage.CalculateTree("feature-a", lineage, configdomain.OrderAsc)
			want := proposallineage.TreeNode{
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
			have := proposallineage.CalculateTree("feature-a", lineage, configdomain.OrderDesc)
			want := proposallineage.TreeNode{
				Branch: "main",
				Children: []proposallineage.TreeNode{
					{
						Branch: "feature-a",
						Children: []proposallineage.TreeNode{
							{
								Branch: "feature-b2",
								Children: []proposallineage.TreeNode{
									{
										Branch:   "feature-b2b",
										Children: []proposallineage.TreeNode{},
									},
									{
										Branch:   "feature-b2a",
										Children: []proposallineage.TreeNode{},
									},
								},
							},
							{
								Branch: "feature-b1",
								Children: []proposallineage.TreeNode{
									{
										Branch:   "feature-b1b",
										Children: []proposallineage.TreeNode{},
									},
									{
										Branch:   "feature-b1a",
										Children: []proposallineage.TreeNode{},
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
			have := proposallineage.CalculateTree("main", lineage, configdomain.OrderAsc)
			want := proposallineage.TreeNode{
				Branch:   "main",
				Children: []proposallineage.TreeNode{},
			}
			must.Eq(t, want, have)
		})
	})
}
