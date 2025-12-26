package proposallineage

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/shoenig/test/must"
)

func TestNewTree2(t *testing.T) {
	t.Parallel()
	t.Run("stand-alone perennial branch", func(t *testing.T) {
		t.Parallel()
		lineage := configdomain.NewLineage()
		have := CalculateTree("main", lineage)
		want := TreeNode2{
			Branch:   "main",
			Children: []TreeNode2{},
		}
		must.Eq(t, want, have)
	})

	// TODO: delete after code is written
	t.Run("branch with one ancestor", func(t *testing.T) {
		t.Parallel()
		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			"feature": "main",
		})
		have := CalculateTree("feature", lineage)
		want := TreeNode2{
			Branch: "main",
			Children: []TreeNode2{
				{
					Branch:   "feature",
					Children: []TreeNode2{},
				},
			},
		}
		must.Eq(t, want, have)
	})

	t.Run("branch with multiple ancestors", func(t *testing.T) {
		t.Parallel()
		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			"feature-a": "main",
			"feature-b": "feature-a",
		})
		have := CalculateTree("feature-b", lineage)
		want := TreeNode2{
			Branch: "main",
			Children: []TreeNode2{
				{
					Branch: "feature-a",
					Children: []TreeNode2{
						{
							Branch:   "feature-b",
							Children: []TreeNode2{},
						},
					},
				},
			},
		}
		must.Eq(t, want, have)
	})

	t.Run("branch with one descendent", func(t *testing.T) {
		t.Parallel()
		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			"feature-a": "main",
			"feature-b": "feature-a",
		})
		have := CalculateTree("feature-a", lineage)
		want := TreeNode2{
			Branch: "main",
			Children: []TreeNode2{
				{
					Branch: "feature-a",
					Children: []TreeNode2{
						{
							Branch:   "feature-b",
							Children: []TreeNode2{},
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
		have := CalculateTree("feature-a", lineage)
		want := TreeNode2{
			Branch: "main",
			Children: []TreeNode2{
				{
					Branch: "feature-a",
					Children: []TreeNode2{
						{
							Branch: "feature-b",
							Children: []TreeNode2{
								{
									Branch: "feature-c",
									Children: []TreeNode2{
										{
											Branch:   "feature-d",
											Children: []TreeNode2{},
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
		have := CalculateTree("feature-a", lineage)
		want := TreeNode2{
			Branch: "main",
			Children: []TreeNode2{
				{
					Branch: "feature-a",
					Children: []TreeNode2{
						{
							Branch: "feature-b1",
							Children: []TreeNode2{
								{
									Branch:   "feature-b1a",
									Children: []TreeNode2{},
								},
								{
									Branch:   "feature-b1b",
									Children: []TreeNode2{},
								},
							},
						},
						{
							Branch: "feature-b2",
							Children: []TreeNode2{
								{
									Branch:   "feature-b2a",
									Children: []TreeNode2{},
								},
								{
									Branch:   "feature-b2b",
									Children: []TreeNode2{},
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
		have := CalculateTree("feature-a", lineage)
		want := TreeNode2{
			Branch: "main",
			Children: []TreeNode2{
				{
					Branch: "feature-a",
					Children: []TreeNode2{
						{
							Branch:   "feature-a1",
							Children: []TreeNode2{},
						},
					},
				},
			},
		}
		must.Eq(t, want, have)
	})
}
