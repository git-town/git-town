package configdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/domain"
	"github.com/shoenig/test/must"
)

func TestLineage(t *testing.T) {
	t.Parallel()
	main := domain.NewLocalBranchName("main")
	one := domain.NewLocalBranchName("one")
	two := domain.NewLocalBranchName("two")
	three := domain.NewLocalBranchName("three")

	t.Run("Ancestors", func(t *testing.T) {
		t.Parallel()
		t.Run("provides all ancestor branches, oldest first", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.Lineage{}
			lineage[three] = two
			lineage[two] = one
			lineage[one] = main
			have := lineage.Ancestors(three)
			want := domain.LocalBranchNames{main, one, two}
			must.Eq(t, want, have)
		})
		t.Run("one ancestor", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.Lineage{}
			lineage[one] = main
			have := lineage.Ancestors(one)
			want := domain.LocalBranchNames{main}
			must.Eq(t, want, have)
		})
		t.Run("no ancestors", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.Lineage{}
			lineage[one] = main
			have := lineage.Ancestors(two)
			want := domain.LocalBranchNames{}
			must.Eq(t, want, have)
		})
	})

	t.Run("BranchAndAncestors", func(t *testing.T) {
		t.Parallel()
		lineage := configdomain.Lineage{}
		lineage[one] = main
		have := lineage.BranchAndAncestors(one)
		want := domain.LocalBranchNames{main, one}
		must.Eq(t, want, have)
	})

	t.Run("BranchesAndAncestors", func(t *testing.T) {
		t.Parallel()
		t.Run("deep lineage, multiple branches", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.Lineage{}
			lineage[one] = main
			lineage[two] = one
			lineage[three] = two
			give := domain.LocalBranchNames{two, one}
			have := lineage.BranchesAndAncestors(give)
			want := domain.LocalBranchNames{main, one, two}
			must.Eq(t, want, have)
		})
		t.Run("deep lineage, multiple branches out of order", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.Lineage{}
			lineage[one] = main
			lineage[two] = one
			lineage[three] = two
			give := domain.LocalBranchNames{one, two, main, three}
			have := lineage.BranchesAndAncestors(give)
			want := domain.LocalBranchNames{main, one, two, three}
			must.Eq(t, want, have)
		})
		t.Run("deep lineage, single branch", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.Lineage{}
			lineage[one] = main
			lineage[two] = one
			lineage[three] = two
			give := domain.LocalBranchNames{three}
			have := lineage.BranchesAndAncestors(give)
			want := domain.LocalBranchNames{main, one, two, three}
			must.Eq(t, want, have)
		})
	})

	t.Run("BranchAndAncestors", func(t *testing.T) {
		t.Parallel()
		lineage := configdomain.Lineage{}
		lineage[one] = main
		have := lineage.BranchAndAncestors(one)
		want := domain.LocalBranchNames{main, one}
		must.Eq(t, want, have)
	})

	t.Run("BranchNames", func(t *testing.T) {
		t.Parallel()
		lineage := configdomain.Lineage{}
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
			lineage := configdomain.Lineage{}
			lineage[twoA] = one
			lineage[twoB] = one
			have := lineage.Children(one)
			want := domain.LocalBranchNames{twoA, twoB}
			must.Eq(t, want, have)
		})
		t.Run("provides only the immediate children, i.e. no grandchildren", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.Lineage{}
			lineage[two] = one
			lineage[three] = two
			have := lineage.Children(one)
			want := domain.LocalBranchNames{two}
			must.Eq(t, want, have)
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.Lineage{}
			have := lineage.Children(one)
			want := domain.LocalBranchNames{}
			must.Eq(t, want, have)
		})
	})

	t.Run("Contains", func(t *testing.T) {
		t.Parallel()
		t.Run("has a parent", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.Lineage{}
			lineage[two] = one
			must.True(t, lineage.HasParents(two))
		})
		t.Run("has no parent", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.Lineage{}
			must.False(t, lineage.HasParents(main))
		})
	})

	t.Run("IsAncestor", func(t *testing.T) {
		t.Run("recognizes greatgrandparent", func(t *testing.T) {
			t.Parallel()
			four := domain.NewLocalBranchName("four")
			lineage := configdomain.Lineage{}
			lineage[four] = three
			lineage[three] = two
			lineage[two] = one
			must.True(t, lineage.IsAncestor(one, four))
		})
		t.Run("child branches are not ancestors", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.Lineage{}
			lineage[two] = one
			must.True(t, lineage.IsAncestor(one, two))
		})
		t.Run("unrelated branches are not ancestors", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.Lineage{}
			lineage[two] = one
			lineage[three] = one
			must.False(t, lineage.IsAncestor(two, three))
		})
	})

	t.Run("OrderedHierarchically", func(t *testing.T) {
		t.Run("multiple lineages", func(t *testing.T) {
			t.Parallel()
			oneA := domain.NewLocalBranchName("oneA")
			oneA1 := domain.NewLocalBranchName("oneA1")
			oneA2 := domain.NewLocalBranchName("oneA2")
			oneB := domain.NewLocalBranchName("oneB")
			lineage := configdomain.Lineage{}
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
		t.Run("deep lineage", func(t *testing.T) {
			t.Parallel()
			one := domain.NewLocalBranchName("one")
			two := domain.NewLocalBranchName("two")
			three := domain.NewLocalBranchName("three")
			four := domain.NewLocalBranchName("four")
			lineage := configdomain.Lineage{}
			lineage[one] = main
			lineage[two] = one
			lineage[three] = two
			lineage[four] = three
			have := domain.LocalBranchNames{four, one}
			want := domain.LocalBranchNames{one, four}
			lineage.OrderHierarchically(have)
			must.Eq(t, want, have)
		})
		t.Run("elements out of order", func(t *testing.T) {
			t.Parallel()
			one := domain.NewLocalBranchName("one")
			two := domain.NewLocalBranchName("two")
			three := domain.NewLocalBranchName("three")
			lineage := configdomain.Lineage{}
			lineage[one] = main
			lineage[two] = one
			lineage[three] = two
			have := domain.LocalBranchNames{one, two, main, three}
			want := domain.LocalBranchNames{main, one, two, three}
			lineage.OrderHierarchically(have)
			must.Eq(t, want, have)
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.Lineage{}
			have := lineage.Parent(main)
			must.True(t, have.IsEmpty())
		})
	})

	t.Run("RemoveBranch", func(t *testing.T) {
		t.Parallel()
		t.Run("branch is a parent branch", func(t *testing.T) {
			t.Parallel()
			main := domain.NewLocalBranchName("main")
			branch1 := domain.NewLocalBranchName("branch-1")
			branch1a := domain.NewLocalBranchName("branch-1a")
			branch1b := domain.NewLocalBranchName("branch-1b")
			branch2 := domain.NewLocalBranchName("branch-2")
			have := configdomain.Lineage{
				branch1:  main,
				branch1a: branch1,
				branch1b: branch1,
				branch2:  main,
			}
			have.RemoveBranch(branch1)
			want := configdomain.Lineage{
				branch1a: main,
				branch1b: main,
				branch2:  main,
			}
			must.Eq(t, want, have)
		})
		t.Run("branch is a child branch", func(t *testing.T) {
			t.Parallel()
			main := domain.NewLocalBranchName("main")
			branch1 := domain.NewLocalBranchName("branch-1")
			branch2 := domain.NewLocalBranchName("branch-2")
			have := configdomain.Lineage{
				branch1: main,
				branch2: main,
			}
			have.RemoveBranch(branch1)
			want := configdomain.Lineage{
				branch2: main,
			}
			must.Eq(t, want, have)
		})
		t.Run("branch is not in lineage", func(t *testing.T) {
			t.Parallel()
			main := domain.NewLocalBranchName("main")
			branch1 := domain.NewLocalBranchName("branch-1")
			branch2 := domain.NewLocalBranchName("branch-2")
			have := configdomain.Lineage{
				branch1: main,
			}
			have.RemoveBranch(branch2)
			want := configdomain.Lineage{
				branch1: main,
			}
			must.Eq(t, want, have)
		})
	})

	t.Run("Roots", func(t *testing.T) {
		t.Parallel()
		t.Run("multiple roots with nested child branches", func(t *testing.T) {
			t.Parallel()
			prod := domain.NewLocalBranchName("prod")
			hotfix1 := domain.NewLocalBranchName("hotfix1")
			hotfix2 := domain.NewLocalBranchName("hotfix2")
			lineage := configdomain.Lineage{}
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
			lineage := configdomain.Lineage{}
			lineage[one] = main
			have := lineage.Roots()
			want := domain.LocalBranchNames{main}
			must.Eq(t, want, have)
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.Lineage{}
			have := lineage.Roots()
			want := domain.LocalBranchNames{}
			must.Eq(t, want, have)
		})
	})
}
