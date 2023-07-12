package config_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/stretchr/testify/assert"
)

func newLineage() config.Lineage {
	return config.Lineage{map[string]string{}, "main"}
}

func TestAncestry(t *testing.T) {
	t.Parallel()

	t.Run("AddParent and Ancestors", func(t *testing.T) {
		t.Parallel()
		t.Run("multiple ancestors", func(t *testing.T) {
			t.Parallel()
			ancestry := newLineage()
			ancestry.SetParent("three", "two")
			ancestry.SetParent("two", "one")
			ancestry.SetParent("one", "main")
			have := ancestry.Ancestors("three")
			want := []string{"main", "one", "two"}
			assert.Equal(t, want, have)
		})
		t.Run("one ancestor", func(t *testing.T) {
			t.Parallel()
			ancestry := newLineage()
			ancestry.SetParent("one", "main")
			have := ancestry.Ancestors("one")
			want := []string{"main"}
			assert.Equal(t, want, have)
		})
		t.Run("no ancestors", func(t *testing.T) {
			t.Parallel()
			ancestry := newLineage()
			ancestry.SetParent("one", "main")
			have := ancestry.Ancestors("two")
			want := []string{}
			assert.Equal(t, want, have)
		})
	})

	t.Run("Children", func(t *testing.T) {
		t.Parallel()
		t.Run("multiple children", func(t *testing.T) {
			t.Parallel()
			ancestry := newLineage()
			ancestry.SetParent("beta1", "alpha")
			ancestry.SetParent("beta2", "alpha")
			have := ancestry.Children("alpha")
			want := []string{"beta1", "beta2"}
			assert.Equal(t, want, have)
		})
		t.Run("child has children", func(t *testing.T) {
			t.Parallel()
			ancestry := newLineage()
			ancestry.SetParent("beta", "alpha")
			ancestry.SetParent("gamma", "beta")
			have := ancestry.Children("alpha")
			want := []string{"beta"}
			assert.Equal(t, want, have)
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			ancestry := newLineage()
			have := ancestry.Children("alpha")
			want := []string{}
			assert.Equal(t, want, have)
		})
	})

	t.Run("HasParent", func(t *testing.T) {
		t.Parallel()
		t.Run("has a parent", func(t *testing.T) {
			t.Parallel()
			ancestry := newLineage()
			ancestry.SetParent("beta", "alpha")
			assert.True(t, ancestry.HasParents("beta"))
		})
		t.Run("has no parent", func(t *testing.T) {
			t.Parallel()
			ancestry := newLineage()
			assert.False(t, ancestry.HasParents("foo"))
		})
	})

	t.Run("IsAncestor", func(t *testing.T) {
		t.Run("greatgrandparent", func(t *testing.T) {
			t.Parallel()
			ancestry := newLineage()
			ancestry.SetParent("four", "three")
			ancestry.SetParent("three", "two")
			ancestry.SetParent("two", "one")
			assert.True(t, ancestry.IsAncestor("four", "one"))
		})
		t.Run("direct parent", func(t *testing.T) {
			t.Parallel()
			ancestry := newLineage()
			ancestry.SetParent("two", "one")
			assert.True(t, ancestry.IsAncestor("two", "one"))
		})
		t.Run("direct child", func(t *testing.T) {
			t.Parallel()
			ancestry := newLineage()
			ancestry.SetParent("child", "parent")
			assert.True(t, ancestry.IsAncestor("one", "two"))
		})
		t.Run("not related", func(t *testing.T) {
			t.Parallel()
			ancestry := newLineage()
			ancestry.SetParent("two", "one")
			ancestry.SetParent("three", "one")
			assert.False(t, ancestry.IsAncestor("three", "two"))
		})
	})

	t.Run("Parent", func(t *testing.T) {
		t.Parallel()
		t.Run("has parent", func(t *testing.T) {
			t.Parallel()
			ancestry := newLineage()
			ancestry.SetParent("two", "one")
			assert.Equal(t, "one", ancestry.Parent("two"))
		})
		t.Run("has no parent", func(t *testing.T) {
			t.Parallel()
			ancestry := newLineage()
			assert.Equal(t, "", ancestry.Parent("foo"))
		})
	})

	t.Run("Roots", func(t *testing.T) {
		t.Parallel()
		t.Run("multiple roots with nested child branches", func(t *testing.T) {
			t.Parallel()
			ancestry := newLineage()
			ancestry.SetParent("two", "one")
			ancestry.SetParent("one", "main")
			ancestry.SetParent("beta", "alpha")
			ancestry.SetParent("alpha", "main")
			ancestry.SetParent("hotfix1", "prod")
			ancestry.SetParent("hotfix2", "prod")
			have := ancestry.Roots()
			want := []string{"main", "prod"}
			assert.Equal(t, want, have)
		})
		t.Run("no nested branches", func(t *testing.T) {
			t.Parallel()
			ancestry := newLineage()
			ancestry.SetParent("one", "main")
			ancestry.SetParent("alpha", "main")
			have := ancestry.Roots()
			want := []string{"main"}
			assert.Equal(t, want, have)
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			ancestry := newLineage()
			have := ancestry.Roots()
			want := []string{}
			assert.Equal(t, want, have)
		})
	})
}
