package proposallineage_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/proposallineage"
	"github.com/shoenig/test/must"
)

func TestCalculateTree(t *testing.T) {
	t.Parallel()

	t.Run("branch with multiple ancestors", func(t *testing.T) {
		t.Parallel()
		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			"feature-a": "main",
			"feature-b": "feature-a",
		})
		have := proposallineage.CalculateTree("feature-b", lineage)
		want := proposallineage.TreeNode2{
			Branch: "main",
			Children: []proposallineage.TreeNode2{
				{
					Branch: "feature-a",
					Children: []proposallineage.TreeNode2{
						{
							Branch:   "feature-b",
							Children: []proposallineage.TreeNode2{},
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
		have := proposallineage.CalculateTree("feature-a", lineage)
		want := proposallineage.TreeNode2{
			Branch: "main",
			Children: []proposallineage.TreeNode2{
				{
					Branch: "feature-a",
					Children: []proposallineage.TreeNode2{
						{
							Branch: "feature-b1",
							Children: []proposallineage.TreeNode2{
								{
									Branch:   "feature-b1a",
									Children: []proposallineage.TreeNode2{},
								},
								{
									Branch:   "feature-b1b",
									Children: []proposallineage.TreeNode2{},
								},
							},
						},
						{
							Branch: "feature-b2",
							Children: []proposallineage.TreeNode2{
								{
									Branch:   "feature-b2a",
									Children: []proposallineage.TreeNode2{},
								},
								{
									Branch:   "feature-b2b",
									Children: []proposallineage.TreeNode2{},
								},
							},
						},
					},
				},
			},
		}
		must.Eq(t, want, have)
	})

	t.Run("branch with multiple descendents", func(t *testing.T) {
		t.Parallel()
		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			"feature-a": "main",
			"feature-b": "feature-a",
			"feature-c": "feature-b",
			"feature-d": "feature-c",
		})
		have := proposallineage.CalculateTree("feature-a", lineage)
		want := proposallineage.TreeNode2{
			Branch: "main",
			Children: []proposallineage.TreeNode2{
				{
					Branch: "feature-a",
					Children: []proposallineage.TreeNode2{
						{
							Branch: "feature-b",
							Children: []proposallineage.TreeNode2{
								{
									Branch: "feature-c",
									Children: []proposallineage.TreeNode2{
										{
											Branch:   "feature-d",
											Children: []proposallineage.TreeNode2{},
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

	t.Run("ignore independent lineages", func(t *testing.T) {
		t.Parallel()
		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			"feature-a":  "main",
			"feature-a1": "feature-a",
			"feature-b":  "main",
			"feature-b1": "feature-b",
		})
		have := proposallineage.CalculateTree("feature-a", lineage)
		want := proposallineage.TreeNode2{
			Branch: "main",
			Children: []proposallineage.TreeNode2{
				{
					Branch: "feature-a",
					Children: []proposallineage.TreeNode2{
						{
							Branch:   "feature-a1",
							Children: []proposallineage.TreeNode2{},
						},
					},
				},
			},
		}
		must.Eq(t, want, have)
	})

	t.Run("stand-alone perennial branch", func(t *testing.T) {
		t.Parallel()
		lineage := configdomain.NewLineage()
		have := proposallineage.CalculateTree("main", lineage)
		want := proposallineage.TreeNode2{
			Branch:   "main",
			Children: []proposallineage.TreeNode2{},
		}
		must.Eq(t, want, have)
	})
}
