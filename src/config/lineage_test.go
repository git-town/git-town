package config_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/stretchr/testify/assert"
)

func TestLineage(t *testing.T) {
	t.Parallel()

	t.Run("Ancestors", func(t *testing.T) {
		t.Parallel()
		t.Run("provides all ancestor branches, oldest first", func(t *testing.T) {
			t.Parallel()
			ancestry := config.Lineage{}
			ancestry["three"] = "two"
			ancestry["two"] = "one"
			ancestry["one"] = "main"
			have := ancestry.Ancestors("three")
			want := []string{"main", "one", "two"}
			assert.Equal(t, want, have)
		})
		t.Run("one ancestor", func(t *testing.T) {
			t.Parallel()
			ancestry := config.Lineage{}
			ancestry["one"] = "main"
			have := ancestry.Ancestors("one")
			want := []string{"main"}
			assert.Equal(t, want, have)
		})
		t.Run("no ancestors", func(t *testing.T) {
			t.Parallel()
			ancestry := config.Lineage{}
			ancestry["one"] = "main"
			have := ancestry.Ancestors("two")
			want := []string{}
			assert.Equal(t, want, have)
		})
	})

	t.Run("Children", func(t *testing.T) {
		t.Parallel()
		t.Run("provides all children of the given branch, ordered alphabetically", func(t *testing.T) {
			t.Parallel()
			ancestry := config.Lineage{}
			ancestry["beta1"] = "alpha"
			ancestry["beta2"] = "alpha"
			have := ancestry.Children("alpha")
			want := []string{"beta1", "beta2"}
			assert.Equal(t, want, have)
		})
		t.Run("provides only the immediate children, i.e. no grandchildren", func(t *testing.T) {
			t.Parallel()
			ancestry := config.Lineage{}
			ancestry["beta"] = "alpha"
			ancestry["gamma"] = "beta"
			have := ancestry.Children("alpha")
			want := []string{"beta"}
			assert.Equal(t, want, have)
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			ancestry := config.Lineage{}
			have := ancestry.Children("alpha")
			want := []string{}
			assert.Equal(t, want, have)
		})
	})

	t.Run("Contains", func(t *testing.T) {
		t.Parallel()
		t.Run("has a parent", func(t *testing.T) {
			t.Parallel()
			ancestry := config.Lineage{}
			ancestry["beta"] = "alpha"
			assert.True(t, ancestry.HasParents("beta"))
		})
		t.Run("has no parent", func(t *testing.T) {
			t.Parallel()
			ancestry := config.Lineage{}
			assert.False(t, ancestry.HasParents("foo"))
		})
	})

	t.Run("IsAncestor", func(t *testing.T) {
		t.Run("recognizes greatgrandparent", func(t *testing.T) {
			t.Parallel()
			ancestry := config.Lineage{}
			ancestry["four"] = "three"
			ancestry["three"] = "two"
			ancestry["two"] = "one"
			assert.True(t, ancestry.IsAncestor("one", "four"))
		})
		t.Run("child branches are not ancestors", func(t *testing.T) {
			t.Parallel()
			ancestry := config.Lineage{}
			ancestry["two"] = "one"
			assert.True(t, ancestry.IsAncestor("one", "two"))
		})
		t.Run("unrelated branches are not ancestors", func(t *testing.T) {
			t.Parallel()
			ancestry := config.Lineage{}
			ancestry["two"] = "one"
			assert.False(t, ancestry.IsAncestor("two", "one"))
		})
		t.Run("not related", func(t *testing.T) {
			t.Parallel()
			ancestry := config.Lineage{}
			ancestry["two"] = "one"
			ancestry["three"] = "one"
			assert.False(t, ancestry.IsAncestor("two", "three"))
		})
	})

	t.Run("OrderedHierarchically", func(t *testing.T) {
		t.Run("complex scenario", func(t *testing.T) {
			t.Parallel()
			lineage := config.Lineage{}
			lineage["main"] = ""
			lineage["1"] = "main"
			lineage["1A"] = "1"
			lineage["1B"] = "1"
			lineage["1A1"] = "1A"
			lineage["1A2"] = "1A"
			lineage["2"] = "main"
			want := []string{"main", "1", "1A", "1B", "1A1", "1A2", "2"}
			have := lineage.OrderedHierarchically().BranchNames()
			assert.Equal(t, want, have)
		})
		t.Run("has no parent", func(t *testing.T) {
			t.Parallel()
			ancestry := config.Lineage{}
			assert.Equal(t, "", ancestry.Parent("foo"))
		})
	})

	t.Run("Roots", func(t *testing.T) {
		t.Parallel()
		t.Run("multiple roots with nested child branches", func(t *testing.T) {
			t.Parallel()
			ancestry := config.Lineage{}
			ancestry["two"] = "one"
			ancestry["one"] = "main"
			ancestry["beta"] = "alpha"
			ancestry["alpha"] = "main"
			ancestry["hotfix1"] = "prod"
			ancestry["hotfix2"] = "prod"
			have := ancestry.Roots()
			want := []string{"main", "prod"}
			assert.Equal(t, want, have)
		})
		t.Run("no nested branches", func(t *testing.T) {
			t.Parallel()
			ancestry := config.Lineage{}
			ancestry["one"] = "main"
			ancestry["alpha"] = "main"
			have := ancestry.Roots()
			want := []string{"main"}
			assert.Equal(t, want, have)
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			ancestry := config.Lineage{}
			have := ancestry.Roots()
			want := []string{}
			assert.Equal(t, want, have)
		})
	})
}
