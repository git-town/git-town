package config_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/stretchr/testify/assert"
)

func TestLineage(t *testing.T) {
	t.Parallel()

	t.Run("AddAncestors", func(t *testing.T) {
		t.Parallel()
		t.Run("single branch", func(t *testing.T) {
			lineage := config.Lineage{}
			lineage["one"] = "main"
			give := []string{"one"}
			have := lineage.AddAncestors(give)
			want := []string{"main", "one"}
			assert.Equal(t, want, have)
		})
	})

	t.Run("Ancestors", func(t *testing.T) {
		t.Parallel()
		t.Run("provides all ancestor branches, oldest first", func(t *testing.T) {
			t.Parallel()
			lineage := config.Lineage{}
			lineage["three"] = "two"
			lineage["two"] = "one"
			lineage["one"] = "main"
			have := lineage.Ancestors("three")
			want := []string{"main", "one", "two"}
			assert.Equal(t, want, have)
		})
		t.Run("one ancestor", func(t *testing.T) {
			t.Parallel()
			lineage := config.Lineage{}
			lineage["one"] = "main"
			have := lineage.Ancestors("one")
			want := []string{"main"}
			assert.Equal(t, want, have)
		})
		t.Run("no ancestors", func(t *testing.T) {
			t.Parallel()
			lineage := config.Lineage{}
			lineage["one"] = "main"
			have := lineage.Ancestors("two")
			want := []string{}
			assert.Equal(t, want, have)
		})
	})

	t.Run("BranchNames", func(t *testing.T) {
		t.Parallel()
		lineage := config.Lineage{}
		lineage["one"] = "1"
		lineage["two"] = "2"
		lineage["three"] = "3"
		have := lineage.BranchNames()
		want := []string{"one", "three", "two"}
		assert.Equal(t, want, have)
	})

	t.Run("Children", func(t *testing.T) {
		t.Parallel()
		t.Run("provides all children of the given branch, ordered alphabetically", func(t *testing.T) {
			t.Parallel()
			lineage := config.Lineage{}
			lineage["beta1"] = "alpha"
			lineage["beta2"] = "alpha"
			have := lineage.Children("alpha")
			want := []string{"beta1", "beta2"}
			assert.Equal(t, want, have)
		})
		t.Run("provides only the immediate children, i.e. no grandchildren", func(t *testing.T) {
			t.Parallel()
			lineage := config.Lineage{}
			lineage["beta"] = "alpha"
			lineage["gamma"] = "beta"
			have := lineage.Children("alpha")
			want := []string{"beta"}
			assert.Equal(t, want, have)
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			lineage := config.Lineage{}
			have := lineage.Children("alpha")
			want := []string{}
			assert.Equal(t, want, have)
		})
	})

	t.Run("Contains", func(t *testing.T) {
		t.Parallel()
		t.Run("has a parent", func(t *testing.T) {
			t.Parallel()
			lineage := config.Lineage{}
			lineage["beta"] = "alpha"
			assert.True(t, lineage.HasParents("beta"))
		})
		t.Run("has no parent", func(t *testing.T) {
			t.Parallel()
			lineage := config.Lineage{}
			assert.False(t, lineage.HasParents("foo"))
		})
	})

	t.Run("IsAncestor", func(t *testing.T) {
		t.Run("recognizes greatgrandparent", func(t *testing.T) {
			t.Parallel()
			lineage := config.Lineage{}
			lineage["four"] = "three"
			lineage["three"] = "two"
			lineage["two"] = "one"
			assert.True(t, lineage.IsAncestor("one", "four"))
		})
		t.Run("child branches are not ancestors", func(t *testing.T) {
			t.Parallel()
			lineage := config.Lineage{}
			lineage["two"] = "one"
			assert.True(t, lineage.IsAncestor("one", "two"))
		})
		t.Run("unrelated branches are not ancestors", func(t *testing.T) {
			t.Parallel()
			lineage := config.Lineage{}
			lineage["two"] = "one"
			lineage["three"] = "one"
			assert.False(t, lineage.IsAncestor("two", "three"))
		})
	})

	t.Run("OrderedHierarchically", func(t *testing.T) {
		t.Run("complex scenario", func(t *testing.T) {
			t.Parallel()
			lineage := config.Lineage{}
			lineage["1"] = "main"
			lineage["1A"] = "1"
			lineage["1B"] = "1"
			lineage["1A1"] = "1A"
			lineage["1A2"] = "1A"
			lineage["2"] = "main"
			want := []string{"1", "1A", "1A1", "1A2", "1B", "2"}
			have := lineage.BranchNames()
			lineage.OrderHierarchically(have)
			assert.Equal(t, want, have)
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			lineage := config.Lineage{}
			assert.Equal(t, "", lineage.Parent("foo"))
		})
	})

	t.Run("Roots", func(t *testing.T) {
		t.Parallel()
		t.Run("multiple roots with nested child branches", func(t *testing.T) {
			t.Parallel()
			lineage := config.Lineage{}
			lineage["two"] = "one"
			lineage["one"] = "main"
			lineage["beta"] = "alpha"
			lineage["alpha"] = "main"
			lineage["hotfix1"] = "prod"
			lineage["hotfix2"] = "prod"
			have := lineage.Roots()
			want := []string{"main", "prod"}
			assert.Equal(t, want, have)
		})
		t.Run("no nested branches", func(t *testing.T) {
			t.Parallel()
			lineage := config.Lineage{}
			lineage["one"] = "main"
			lineage["alpha"] = "main"
			have := lineage.Roots()
			want := []string{"main"}
			assert.Equal(t, want, have)
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			lineage := config.Lineage{}
			have := lineage.Roots()
			want := []string{}
			assert.Equal(t, want, have)
		})
	})
}
