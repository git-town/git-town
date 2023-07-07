package config_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/stretchr/testify/assert"
)

func TestAncestry(t *testing.T) {
	t.Parallel()

	t.Run("AddParent and Ancestors", func(t *testing.T) {
		t.Parallel()
		t.Run("multiple ancestors", func(t *testing.T) {
			t.Parallel()
			ancestry := config.NewLineage("main")
			ancestry.AddParent("three", "two")
			ancestry.AddParent("two", "one")
			ancestry.AddParent("one", "main")
			have := ancestry.Ancestors("three")
			want := []string{"main", "one", "two"}
			assert.Equal(t, want, have)
		})
		t.Run("one ancestor", func(t *testing.T) {
			t.Parallel()
			ancestry := config.NewLineage("main")
			ancestry.AddParent("one", "main")
			have := ancestry.Ancestors("one")
			want := []string{"main"}
			assert.Equal(t, want, have)
		})
		t.Run("no ancestors", func(t *testing.T) {
			t.Parallel()
			ancestry := config.NewLineage("main")
			ancestry.AddParent("one", "main")
			have := ancestry.Ancestors("two")
			want := []string{}
			assert.Equal(t, want, have)
		})
	})

	t.Run("Roots", func(t *testing.T) {
		t.Parallel()
		t.Run("multiple roots with nested child branches", func(t *testing.T) {
			t.Parallel()
			ancestry := config.NewLineage("main")
			ancestry.AddParent("two", "one")
			ancestry.AddParent("one", "main")
			ancestry.AddParent("beta", "alpha")
			ancestry.AddParent("alpha", "main")
			ancestry.AddParent("hotfix1", "prod")
			ancestry.AddParent("hotfix2", "prod")
			have := ancestry.Roots()
			want := []string{"main", "prod"}
			assert.Equal(t, want, have)
		})
		t.Run("no nested branches", func(t *testing.T) {
			t.Parallel()
			ancestry := config.NewLineage("main")
			ancestry.AddParent("one", "main")
			ancestry.AddParent("alpha", "main")
			have := ancestry.Roots()
			want := []string{"main"}
			assert.Equal(t, want, have)
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			ancestry := config.NewLineage("main")
			have := ancestry.Roots()
			want := []string{}
			assert.Equal(t, want, have)
		})
	})

	t.Run("Children", func(t *testing.T) {
		t.Parallel()
		t.Run("multiple children", func(t *testing.T) {
			t.Parallel()
			ancestry := config.NewLineage("main")
			ancestry.AddParent("beta1", "alpha")
			ancestry.AddParent("beta2", "alpha")
			have := ancestry.Children("alpha")
			want := []string{"beta1", "beta2"}
			assert.Equal(t, want, have)
		})
		t.Run("child has children", func(t *testing.T) {
			t.Parallel()
			ancestry := config.NewLineage("main")
			ancestry.AddParent("beta", "alpha")
			ancestry.AddParent("gamma", "beta")
			have := ancestry.Children("alpha")
			want := []string{"beta"}
			assert.Equal(t, want, have)
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			ancestry := config.NewLineage("main")
			have := ancestry.Children("alpha")
			want := []string{}
			assert.Equal(t, want, have)
		})
	})

	t.Run("HasParent", func(t *testing.T) {
		t.Parallel()
		t.Run("has a parent", func(t *testing.T) {
			t.Parallel()
			ancestry := config.NewLineage("main")
			ancestry.AddParent("beta", "alpha")
			assert.True(t, ancestry.HasParents("beta"))
		})
		t.Run("has no parent", func(t *testing.T) {
			t.Parallel()
			ancestry := config.NewLineage("main")
			assert.False(t, ancestry.HasParents("foo"))
		})
	})

	t.Run("IsAncestor", func(t *testing.T) {
		t.Run("greatgrandparent", func(t *testing.T) {
			t.Parallel()
			ancestry := config.NewLineage("main")
			ancestry.AddParent("four", "three")
			ancestry.AddParent("three", "two")
			ancestry.AddParent("two", "one")
			assert.True(t, ancestry.IsAncestor("four", "one"))
		})
		t.Run("direct parent", func(t *testing.T) {
			t.Parallel()
			ancestry := config.NewLineage("main")
			ancestry.AddParent("two", "one")
			assert.True(t, ancestry.IsAncestor("two", "one"))
		})
	})

	t.Run("Parent", func(t *testing.T) {
		t.Parallel()
		t.Run("has parent", func(t *testing.T) {
			t.Parallel()
			ancestry := config.NewLineage("main")
			ancestry.AddParent("two", "one")
			assert.Equal(t, "one", ancestry.Parent("two"))
		})
		t.Run("has no parent", func(t *testing.T) {
			t.Parallel()
			ancestry := config.NewLineage("main")
			assert.Equal(t, "", ancestry.Parent("foo"))
		})
	})
}
