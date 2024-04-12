package configdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v13/src/config/configdomain"
	"github.com/git-town/git-town/v13/src/git/gitdomain"
	"github.com/shoenig/test/must"
)

func TestLineage(t *testing.T) {
	t.Parallel()
	main := gitdomain.NewLocalBranchName("main")
	one := gitdomain.NewLocalBranchName("one")
	two := gitdomain.NewLocalBranchName("two")
	three := gitdomain.NewLocalBranchName("three")

	t.Run("Ancestors", func(t *testing.T) {
		t.Parallel()
		t.Run("provides all ancestor branches, oldest first", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.Lineage{}
			lineage[three] = two
			lineage[two] = one
			lineage[one] = main
			have := lineage.Ancestors(three)
			want := gitdomain.LocalBranchNames{main, one, two}
			must.Eq(t, want, have)
		})
		t.Run("one ancestor", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.Lineage{}
			lineage[one] = main
			have := lineage.Ancestors(one)
			want := gitdomain.LocalBranchNames{main}
			must.Eq(t, want, have)
		})
		t.Run("no ancestors", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.Lineage{}
			lineage[one] = main
			have := lineage.Ancestors(two)
			want := gitdomain.LocalBranchNames{}
			must.Eq(t, want, have)
		})
	})

	t.Run("AncestorsWithoutRoot", func(t *testing.T) {
		t.Parallel()
		t.Run("provides all ancestor branches, oldest first", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.Lineage{}
			lineage[three] = two
			lineage[two] = one
			lineage[one] = main
			have := lineage.AncestorsWithoutRoot(three)
			want := gitdomain.LocalBranchNames{one, two}
			must.Eq(t, want, have)
		})
		t.Run("one ancestor", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.Lineage{}
			lineage[one] = main
			have := lineage.AncestorsWithoutRoot(one)
			want := gitdomain.LocalBranchNames{}
			must.Eq(t, want, have)
		})
		t.Run("no ancestors", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.Lineage{}
			lineage[one] = main
			have := lineage.AncestorsWithoutRoot(two)
			want := gitdomain.LocalBranchNames{}
			must.Eq(t, want, have)
		})
	})

	t.Run("BranchAndAncestors", func(t *testing.T) {
		t.Parallel()
		lineage := configdomain.Lineage{}
		lineage[one] = main
		have := lineage.BranchAndAncestors(one)
		want := gitdomain.LocalBranchNames{main, one}
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
			give := gitdomain.LocalBranchNames{two, one}
			have := lineage.BranchesAndAncestors(give)
			want := gitdomain.LocalBranchNames{main, one, two}
			must.Eq(t, want, have)
		})
		t.Run("deep lineage, multiple branches out of order", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.Lineage{}
			lineage[one] = main
			lineage[two] = one
			lineage[three] = two
			give := gitdomain.LocalBranchNames{one, two, main, three}
			have := lineage.BranchesAndAncestors(give)
			want := gitdomain.LocalBranchNames{main, one, two, three}
			must.Eq(t, want, have)
		})
		t.Run("deep lineage, single branch", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.Lineage{}
			lineage[one] = main
			lineage[two] = one
			lineage[three] = two
			give := gitdomain.LocalBranchNames{three}
			have := lineage.BranchesAndAncestors(give)
			want := gitdomain.LocalBranchNames{main, one, two, three}
			must.Eq(t, want, have)
		})
	})

	t.Run("BranchAndAncestors", func(t *testing.T) {
		t.Parallel()
		lineage := configdomain.Lineage{}
		lineage[one] = main
		have := lineage.BranchAndAncestors(one)
		want := gitdomain.LocalBranchNames{main, one}
		must.Eq(t, want, have)
	})

	t.Run("BranchNames", func(t *testing.T) {
		t.Parallel()
		lineage := configdomain.Lineage{}
		lineage[one] = main
		lineage[two] = main
		lineage[three] = main
		have := lineage.BranchNames()
		want := gitdomain.LocalBranchNames{one, three, two}
		must.Eq(t, want, have)
	})

	t.Run("Children", func(t *testing.T) {
		t.Parallel()
		t.Run("provides all children of the given branch, ordered alphabetically", func(t *testing.T) {
			t.Parallel()
			twoA := gitdomain.NewLocalBranchName("twoA")
			twoB := gitdomain.NewLocalBranchName("twoB")
			lineage := configdomain.Lineage{}
			lineage[twoA] = one
			lineage[twoB] = one
			have := lineage.Children(one)
			want := gitdomain.LocalBranchNames{twoA, twoB}
			must.Eq(t, want, have)
		})
		t.Run("provides only the immediate children, i.e. no grandchildren", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.Lineage{}
			lineage[two] = one
			lineage[three] = two
			have := lineage.Children(one)
			want := gitdomain.LocalBranchNames{two}
			must.Eq(t, want, have)
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.Lineage{}
			have := lineage.Children(one)
			want := gitdomain.LocalBranchNames{}
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

	t.Run("Descendants", func(t *testing.T) {
		t.Parallel()
		t.Run("branch has no children", func(t *testing.T) {
			t.Parallel()
			branch := gitdomain.NewLocalBranchName("branch")
			other := gitdomain.NewLocalBranchName("other")
			lineage := configdomain.Lineage{
				branch: main,
				other:  main,
			}
			have := lineage.Descendants(branch)
			want := gitdomain.LocalBranchNames{}
			must.Eq(t, want, have)
		})
		t.Run("branch has only direct children", func(t *testing.T) {
			t.Parallel()
			branch := gitdomain.NewLocalBranchName("branch")
			child1 := gitdomain.NewLocalBranchName("child1")
			child2 := gitdomain.NewLocalBranchName("child2")
			other := gitdomain.NewLocalBranchName("other")
			lineage := configdomain.Lineage{
				branch: main,
				child1: branch,
				child2: branch,
				other:  main,
			}
			have := lineage.Descendants(branch)
			want := gitdomain.LocalBranchNames{child1, child2}
			must.Eq(t, want, have)
		})
		t.Run("branch has grandchildren", func(t *testing.T) {
			t.Parallel()
			branch := gitdomain.NewLocalBranchName("branch")
			child1 := gitdomain.NewLocalBranchName("child1")
			child1a := gitdomain.NewLocalBranchName("child1a")
			child1b := gitdomain.NewLocalBranchName("child1b")
			child2 := gitdomain.NewLocalBranchName("child2")
			child2a := gitdomain.NewLocalBranchName("child2a")
			child2b := gitdomain.NewLocalBranchName("child2b")
			other := gitdomain.NewLocalBranchName("other")
			lineage := configdomain.Lineage{
				branch:  main,
				child1:  branch,
				child1a: child1,
				child1b: child1,
				child2:  branch,
				child2a: child2,
				child2b: child2,
				other:   main,
			}
			have := lineage.Descendants(branch)
			want := gitdomain.LocalBranchNames{child1, child1a, child1b, child2, child2a, child2b}
			must.Eq(t, want, have)
		})
	})

	t.Run("IsAncestor", func(t *testing.T) {
		t.Run("recognizes greatgrandparent", func(t *testing.T) {
			t.Parallel()
			four := gitdomain.NewLocalBranchName("four")
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
			oneA := gitdomain.NewLocalBranchName("oneA")
			oneA1 := gitdomain.NewLocalBranchName("oneA1")
			oneA2 := gitdomain.NewLocalBranchName("oneA2")
			oneB := gitdomain.NewLocalBranchName("oneB")
			lineage := configdomain.Lineage{}
			lineage[one] = main
			lineage[oneA] = one
			lineage[oneB] = one
			lineage[oneA1] = oneA
			lineage[oneA2] = oneA
			lineage[two] = main
			have := lineage.BranchNames()
			want := gitdomain.LocalBranchNames{one, oneA, oneA1, oneA2, oneB, two}
			lineage.OrderHierarchically(have)
			must.Eq(t, want, have)
		})
		t.Run("deep lineage", func(t *testing.T) {
			t.Parallel()
			one := gitdomain.NewLocalBranchName("one")
			two := gitdomain.NewLocalBranchName("two")
			three := gitdomain.NewLocalBranchName("three")
			four := gitdomain.NewLocalBranchName("four")
			lineage := configdomain.Lineage{}
			lineage[one] = main
			lineage[two] = one
			lineage[three] = two
			lineage[four] = three
			have := gitdomain.LocalBranchNames{four, one}
			want := gitdomain.LocalBranchNames{one, four}
			lineage.OrderHierarchically(have)
			must.Eq(t, want, have)
		})
		t.Run("elements out of order", func(t *testing.T) {
			t.Parallel()
			one := gitdomain.NewLocalBranchName("one")
			two := gitdomain.NewLocalBranchName("two")
			three := gitdomain.NewLocalBranchName("three")
			lineage := configdomain.Lineage{}
			lineage[one] = main
			lineage[two] = one
			lineage[three] = two
			have := gitdomain.LocalBranchNames{one, two, main, three}
			want := gitdomain.LocalBranchNames{main, one, two, three}
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
			main := gitdomain.NewLocalBranchName("main")
			branch1 := gitdomain.NewLocalBranchName("branch-1")
			branch1a := gitdomain.NewLocalBranchName("branch-1a")
			branch1b := gitdomain.NewLocalBranchName("branch-1b")
			branch2 := gitdomain.NewLocalBranchName("branch-2")
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
			main := gitdomain.NewLocalBranchName("main")
			branch1 := gitdomain.NewLocalBranchName("branch-1")
			branch2 := gitdomain.NewLocalBranchName("branch-2")
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
			main := gitdomain.NewLocalBranchName("main")
			branch1 := gitdomain.NewLocalBranchName("branch-1")
			branch2 := gitdomain.NewLocalBranchName("branch-2")
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
		t.Run("multiple roots with child branches", func(t *testing.T) {
			t.Parallel()
			prod := gitdomain.NewLocalBranchName("prod")
			hotfix1 := gitdomain.NewLocalBranchName("hotfix1")
			hotfix2 := gitdomain.NewLocalBranchName("hotfix2")
			lineage := configdomain.Lineage{}
			lineage[two] = one
			lineage[one] = main
			lineage[hotfix1] = prod
			lineage[hotfix2] = prod
			have := lineage.Roots()
			want := gitdomain.LocalBranchNames{main, prod}
			must.Eq(t, want, have)
		})
		t.Run("no stacked changes", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.Lineage{}
			lineage[one] = main
			have := lineage.Roots()
			want := gitdomain.LocalBranchNames{main}
			must.Eq(t, want, have)
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.Lineage{}
			have := lineage.Roots()
			want := gitdomain.LocalBranchNames{}
			must.Eq(t, want, have)
		})
	})
}
