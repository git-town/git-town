package config_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/stretchr/testify/assert"
)

func TestLineage(t *testing.T) {

	t.Run("Ancestors", func(t *testing.T) {
		t.Parallel()
		t.Run("provides all ancestor branches, oldest first", func(t *testing.T) {
			t.Parallel()
			lineage := config.Lineage{
				config.BranchWithParent{
					Name:   "three",
					Parent: "two",
				},
				config.BranchWithParent{
					Name:   "two",
					Parent: "one",
				},
				config.BranchWithParent{
					Name:   "one",
					Parent: "main",
				},
				config.BranchWithParent{
					Name:   "main",
					Parent: "",
				},
				config.BranchWithParent{
					Name:   "other",
					Parent: "one",
				},
			}
			want := []string{"main", "one", "two", "three"}
			have := lineage.Ancestors("three").BranchNames()
			assert.Equal(t, want, have)
		})
		t.Run("no ancestors", func(t *testing.T) {
			t.Parallel()
			lineage := config.Lineage{}
			have := lineage.Ancestors("foo")
			want := config.Lineage{}
			assert.Equal(t, want, have)
		})
	})

	t.Run("BranchNames", func(t *testing.T) {
		t.Parallel()
		t.Run("returns the names of all branches in this collection, ordered the same as the collection", func(t *testing.T) {
			t.Parallel()
			lineage := config.Lineage{
				config.BranchWithParent{Name: "one"},
				config.BranchWithParent{Name: "two"},
				config.BranchWithParent{Name: "three"},
			}
			want := []string{"one", "two", "three"}
			have := lineage.BranchNames()
			assert.Equal(t, want, have)
		})
	})

	t.Run("Children", func(t *testing.T) {
		t.Parallel()
		t.Run("provides all children of the given branch, ordered alphabetically", func(t *testing.T) {
			t.Parallel()
			lineage := config.Lineage{
				config.BranchWithParent{
					Name:   "alpha",
					Parent: "",
				},
				config.BranchWithParent{
					Name:   "beta2",
					Parent: "alpha",
				},
				config.BranchWithParent{
					Name:   "beta1",
					Parent: "alpha",
				},
			}
			have := lineage.Children("alpha").BranchNames()
			want := []string{"beta1", "beta2"}
			assert.Equal(t, want, have)
		})
		t.Run("provides only the immediate children, i.e. no grandchildren", func(t *testing.T) {
			t.Parallel()
			lineage := config.Lineage{
				config.BranchWithParent{
					Name:   "one",
					Parent: "",
				},
				config.BranchWithParent{
					Name:   "two",
					Parent: "one",
				},
				config.BranchWithParent{
					Name:   "three",
					Parent: "two",
				},
			}
			have := lineage.Children("one").BranchNames()
			want := []string{"two"}
			assert.Equal(t, want, have)
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			lineage := config.Lineage{}
			have := lineage.Children("alpha")
			want := config.Lineage{}
			assert.Equal(t, want, have)
		})
	})

	t.Run("Contains", func(t *testing.T) {
		t.Parallel()
		lineage := config.Lineage{
			config.BranchWithParent{Name: "one"},
			config.BranchWithParent{Name: "two"},
		}
		assert.True(t, lineage.Contains("one"))
		assert.True(t, lineage.Contains("two"))
		assert.False(t, lineage.Contains("zonk"))
	})

	t.Run("IsAncestor", func(t *testing.T) {
		t.Run("recognizes greatgrandparent", func(t *testing.T) {
			t.Parallel()
			lineage := config.Lineage{
				config.BranchWithParent{
					Name:   "one",
					Parent: "",
				},
				config.BranchWithParent{
					Name:   "two",
					Parent: "one",
				},
				config.BranchWithParent{
					Name:   "three",
					Parent: "two",
				},
				config.BranchWithParent{
					Name:   "four",
					Parent: "three",
				},
			}
			assert.True(t, lineage.IsAncestor("one", "four"))
		})
		t.Run("child branches are not ancestors", func(t *testing.T) {
			t.Parallel()
			lineage := config.Lineage{
				config.BranchWithParent{
					Name:   "one",
					Parent: "",
				},
				config.BranchWithParent{
					Name:   "two",
					Parent: "one",
				},
			}
			assert.False(t, lineage.IsAncestor("two", "one"))
		})
		t.Run("unrelated branches are not ancestors", func(t *testing.T) {
			t.Parallel()
			lineage := config.Lineage{
				config.BranchWithParent{
					Name:   "one",
					Parent: "",
				},
				config.BranchWithParent{
					Name:   "two",
					Parent: "one",
				},
				config.BranchWithParent{
					Name:   "two",
					Parent: "one",
				},
			}
			assert.False(t, lineage.IsAncestor("two", "three"))
		})
	})

	t.Run("OrderedHierarchically", func(t *testing.T) {
		t.Parallel()
		lineage := config.Lineage{
			config.BranchWithParent{
				Name:   "main",
				Parent: "",
			},
			config.BranchWithParent{
				Name:   "1",
				Parent: "main",
			},
			config.BranchWithParent{
				Name:   "1A",
				Parent: "1",
			},
			config.BranchWithParent{
				Name:   "1B",
				Parent: "one",
			},
			config.BranchWithParent{
				Name:   "1A1",
				Parent: "1A",
			},
			config.BranchWithParent{
				Name:   "1A2",
				Parent: "1A",
			},
			config.BranchWithParent{
				Name:   "2",
				Parent: "main",
			},
		}
		want := []string{"main", "1", "2", "1A", "1B", "1A1", "1A2"}
		have := lineage.OrderedHierarchically().BranchNames()
		assert.Equal(t, want, have)
	})

	t.Run("Roots", func(t *testing.T) {
		t.Parallel()
		t.Run("multiple roots with nested child branches", func(t *testing.T) {
			t.Parallel()
			lineage := config.Lineage{
				config.BranchWithParent{
					Name:   "main",
					Parent: "",
				},
				config.BranchWithParent{
					Name:   "one",
					Parent: "main",
				},
				config.BranchWithParent{
					Name:   "two",
					Parent: "one",
				},
				config.BranchWithParent{
					Name:   "alpha",
					Parent: "main",
				},
				config.BranchWithParent{
					Name:   "beta",
					Parent: "alpha",
				},
				config.BranchWithParent{
					Name:   "prod",
					Parent: "",
				},
				config.BranchWithParent{
					Name:   "hotfix1",
					Parent: "prod",
				},
			}
			want := []string{"main", "prod"}
			have := lineage.Roots().BranchNames()
			assert.Equal(t, want, have)
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			lineage := config.Lineage{}
			have := lineage.Roots()
			want := config.Lineage{}
			assert.Equal(t, want, have)
		})
	})
}
