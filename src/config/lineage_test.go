package config_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/shoenig/test/must"
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
			give := domain.LocalBranchNames{two, one}
			have := lineage.BranchesAndAncestors(give)
			want := domain.LocalBranchNames{main, one, two}
			must.Eq(t, want, have)
		})
	})

	t.Run("BranchAndAncestors", func(t *testing.T) {
		t.Parallel()
		lineage := config.Lineage{}
		lineage[one] = main
		have := lineage.BranchAndAncestors(one)
		want := domain.LocalBranchNames{main, one}
		must.Eq(t, want, have)
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
			want := domain.LocalBranchNames{main, one, two}
			must.Eq(t, want, have)
		})
		t.Run("one ancestor", func(t *testing.T) {
			t.Parallel()
			lineage := config.Lineage{}
			lineage[one] = main
			have := lineage.Ancestors(one)
			want := domain.LocalBranchNames{main}
			must.Eq(t, want, have)
		})
		t.Run("no ancestors", func(t *testing.T) {
			t.Parallel()
			lineage := config.Lineage{}
			lineage[one] = main
			have := lineage.Ancestors(two)
			want := domain.LocalBranchNames{}
			must.Eq(t, want, have)
		})
	})

	t.Run("BranchNames", func(t *testing.T) {
		t.Parallel()
		lineage := config.Lineage{}
		lineage[one] = main
		lineage[two] = main
		lineage[three] = main
		have := lineage.BranchNames()
		want := domain.LocalBranchNames{one, three, two}
		must.Eq(t, want, have)
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
			want := domain.LocalBranchNames{twoA, twoB}
			must.Eq(t, want, have)
		})
		t.Run("provides only the immediate children, i.e. no grandchildren", func(t *testing.T) {
			t.Parallel()
			lineage := config.Lineage{}
			lineage[two] = one
			lineage[three] = two
			have := lineage.Children(one)
			want := domain.LocalBranchNames{two}
			must.Eq(t, want, have)
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			lineage := config.Lineage{}
			have := lineage.Children(one)
			want := domain.LocalBranchNames{}
			must.Eq(t, want, have)
		})
	})

	t.Run("Contains", func(t *testing.T) {
		t.Parallel()
		t.Run("has a parent", func(t *testing.T) {
			t.Parallel()
			lineage := config.Lineage{}
			lineage[two] = one
			must.True(t, lineage.HasParents(two))
		})
		t.Run("has no parent", func(t *testing.T) {
			t.Parallel()
			lineage := config.Lineage{}
			must.False(t, lineage.HasParents(main))
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
			must.True(t, lineage.IsAncestor(one, four))
		})
		t.Run("child branches are not ancestors", func(t *testing.T) {
			t.Parallel()
			lineage := config.Lineage{}
			lineage[two] = one
			must.True(t, lineage.IsAncestor(one, two))
		})
		t.Run("unrelated branches are not ancestors", func(t *testing.T) {
			t.Parallel()
			lineage := config.Lineage{}
			lineage[two] = one
			lineage[three] = one
			must.False(t, lineage.IsAncestor(two, three))
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
			have := lineage.BranchNames()
			want := domain.LocalBranchNames{one, oneA, oneA1, oneA2, oneB, two}
			lineage.OrderHierarchically(have)
			must.Eq(t, want, have)
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			lineage := config.Lineage{}
			have := lineage.Parent(main)
			must.True(t, have.IsEmpty())
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
			want := domain.LocalBranchNames{main, prod}
			must.Eq(t, want, have)
		})
		t.Run("no nested branches", func(t *testing.T) {
			t.Parallel()
			lineage := config.Lineage{}
			lineage[one] = main
			have := lineage.Roots()
			want := domain.LocalBranchNames{main}
			must.Eq(t, want, have)
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			lineage := config.Lineage{}
			have := lineage.Roots()
			want := domain.LocalBranchNames{}
			must.Eq(t, want, have)
		})
	})
}
