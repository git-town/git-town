package config_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/stretchr/testify/assert"
)

func TestLineage(t *testing.T) {
	t.Parallel()
	main := domain.NewLocalBranchName("main")
	one := domain.NewLocalBranchName("one")
	two := domain.NewLocalBranchName("two")
	three := domain.NewLocalBranchName("three")

	t.Run("BranchesAndAncestors", func(t *testing.T) {
		t.Parallel()
		t.Run("many branches", func(t *testing.T) {
			t.Parallel()
			lineage := config.Lineage{}
			lineage[one] = main
			lineage[two] = one
			give := domain.NewLocalBranchNames("two", "one")
			have := lineage.BranchesAndAncestors(give)
			want := domain.NewLocalBranchNames("main", "one", "two")
			assert.Equal(t, want, have)
		})
	})

	t.Run("BranchAndAncestors", func(t *testing.T) {
		t.Parallel()
		lineage := config.Lineage{}
		lineage[one] = main
		have := lineage.BranchAndAncestors(one)
		want := domain.NewLocalBranchNames("main", "one")
		assert.Equal(t, want, have)
	})

	t.Run("Ancestors", func(t *testing.T) {
		t.Parallel()
		t.Run("provides all ancestor branches, oldest first", func(t *testing.T) {
			t.Parallel()
			lineage := config.Lineage{}
			lineage[three] = two
			lineage[two] = one
			lineage[one] = main
			have := lineage.Ancestors(three)
			want := domain.NewLocalBranchNames("main", "one", "two")
			assert.Equal(t, want, have)
		})
		t.Run("one ancestor", func(t *testing.T) {
			t.Parallel()
			lineage := config.Lineage{}
			lineage[one] = main
			have := lineage.Ancestors(one)
			want := domain.NewLocalBranchNames("main")
			assert.Equal(t, want, have)
		})
		t.Run("no ancestors", func(t *testing.T) {
			t.Parallel()
			lineage := config.Lineage{}
			lineage[one] = main
			have := lineage.Ancestors(two)
			want := domain.LocalBranchNames{}
			assert.Equal(t, want, have)
		})
	})

	t.Run("BranchNames", func(t *testing.T) {
		t.Parallel()
		lineage := config.Lineage{}
		lineage[one] = main
		lineage[two] = main
		lineage[three] = main
		have := lineage.BranchNames()
		want := domain.NewLocalBranchNames("one", "three", "two")
		assert.Equal(t, want, have)
	})

	t.Run("Children", func(t *testing.T) {
		t.Parallel()
		t.Run("provides all children of the given branch, ordered alphabetically", func(t *testing.T) {
			t.Parallel()
			twoA := domain.NewLocalBranchName("twoA")
			twoB := domain.NewLocalBranchName("twoB")
			lineage := config.Lineage{}
			lineage[twoA] = one
			lineage[twoB] = one
			have := lineage.Children(one)
			want := domain.NewLocalBranchNames("twoA", "twoB")
			assert.Equal(t, want, have)
		})
		t.Run("provides only the immediate children, i.e. no grandchildren", func(t *testing.T) {
			t.Parallel()
			lineage := config.Lineage{}
			lineage[two] = one
			lineage[three] = two
			have := lineage.Children(one)
			want := domain.NewLocalBranchNames("two")
			assert.Equal(t, want, have)
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			lineage := config.Lineage{}
			have := lineage.Children(one)
			want := domain.LocalBranchNames{}
			assert.Equal(t, want, have)
		})
	})

	t.Run("Contains", func(t *testing.T) {
		t.Parallel()
		t.Run("has a parent", func(t *testing.T) {
			t.Parallel()
			lineage := config.Lineage{}
			lineage[two] = one
			assert.True(t, lineage.HasParents(two))
		})
		t.Run("has no parent", func(t *testing.T) {
			t.Parallel()
			lineage := config.Lineage{}
			assert.False(t, lineage.HasParents(main))
		})
	})

	t.Run("IsAncestor", func(t *testing.T) {
		t.Run("recognizes greatgrandparent", func(t *testing.T) {
			t.Parallel()
			four := domain.NewLocalBranchName("four")
			lineage := config.Lineage{}
			lineage[four] = three
			lineage[three] = two
			lineage[two] = one
			assert.True(t, lineage.IsAncestor(one, four))
		})
		t.Run("child branches are not ancestors", func(t *testing.T) {
			t.Parallel()
			lineage := config.Lineage{}
			lineage[two] = one
			assert.True(t, lineage.IsAncestor(one, two))
		})
		t.Run("unrelated branches are not ancestors", func(t *testing.T) {
			t.Parallel()
			lineage := config.Lineage{}
			lineage[two] = one
			lineage[three] = one
			assert.False(t, lineage.IsAncestor(two, three))
		})
	})

	t.Run("OrderedHierarchically", func(t *testing.T) {
		t.Run("complex scenario", func(t *testing.T) {
			t.Parallel()
			oneA := domain.NewLocalBranchName("oneA")
			oneA1 := domain.NewLocalBranchName("oneA1")
			oneA2 := domain.NewLocalBranchName("oneA2")
			oneB := domain.NewLocalBranchName("oneB")
			lineage := config.Lineage{}
			lineage[one] = main
			lineage[oneA] = one
			lineage[oneB] = one
			lineage[oneA1] = oneA
			lineage[oneA2] = oneA
			lineage[two] = main
			want := domain.LocalBranchNames{one, oneA, oneA1, oneA2, oneB, two}
			have := lineage.BranchNames()
			lineage.OrderHierarchically(have)
			assert.Equal(t, want, have)
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			lineage := config.Lineage{}
			have := lineage.Parent(main)
			assert.True(t, have.IsEmpty())
		})
	})

	t.Run("Roots", func(t *testing.T) {
		t.Parallel()
		t.Run("multiple roots with nested child branches", func(t *testing.T) {
			t.Parallel()
			prod := domain.NewLocalBranchName("prod")
			hotfix1 := domain.NewLocalBranchName("hotfix1")
			hotfix2 := domain.NewLocalBranchName("hotfix2")
			lineage := config.Lineage{}
			lineage[two] = one
			lineage[one] = main
			lineage[hotfix1] = prod
			lineage[hotfix2] = prod
			have := lineage.Roots()
			want := domain.NewLocalBranchNames("main", "prod")
			assert.Equal(t, want, have)
		})
		t.Run("no nested branches", func(t *testing.T) {
			t.Parallel()
			lineage := config.Lineage{}
			lineage[one] = main
			have := lineage.Roots()
			want := domain.NewLocalBranchNames("main")
			assert.Equal(t, want, have)
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			lineage := config.Lineage{}
			have := lineage.Roots()
			want := domain.LocalBranchNames{}
			assert.Equal(t, want, have)
		})
	})
}
