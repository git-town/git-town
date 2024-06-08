package configdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/shoenig/test/must"
)

func TestLineage(t *testing.T) {
	t.Parallel()
	main := gitdomain.NewLocalBranchName("main")
	one := gitdomain.NewLocalBranchName("one")
	two := gitdomain.NewLocalBranchName("two")
	three := gitdomain.NewLocalBranchName("three")

	t.Run("AddParent", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.NewLineage()
			branch := gitdomain.NewLocalBranchName("branch")
			parent := gitdomain.NewLocalBranchName("parent")
			lineage.AddParent(branch, parent)
			have, has := lineage.Parent(branch).Get()
			must.True(t, has)
			must.Eq(t, parent, have)
		})
		t.Run("entry already exists", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.NewLineage()
			branch := gitdomain.NewLocalBranchName("branch")
			parent := gitdomain.NewLocalBranchName("parent")
			lineage.AddParent(branch, parent)
			lineage.AddParent(branch, parent)
			must.EqOp(t, 1, lineage.Len())
			have, has := lineage.Parent(branch).Get()
			must.True(t, has)
			must.Eq(t, parent, have)
		})
	})

	t.Run("Ancestors", func(t *testing.T) {
		t.Parallel()
		t.Run("provides all ancestor branches, oldest first", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.NewLineage()
			lineage.AddParent(three, two)
			lineage.AddParent(two, one)
			lineage.AddParent(one, main)
			have := lineage.Ancestors(three)
			want := gitdomain.LocalBranchNames{main, one, two}
			must.Eq(t, want, have)
		})
		t.Run("one ancestor", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.NewLineage()
			lineage.AddParent(one, main)
			have := lineage.Ancestors(one)
			want := gitdomain.LocalBranchNames{main}
			must.Eq(t, want, have)
		})
		t.Run("no ancestors", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.NewLineage()
			lineage.AddParent(one, main)
			have := lineage.Ancestors(two)
			want := gitdomain.LocalBranchNames{}
			must.Eq(t, want, have)
		})
	})

	t.Run("AncestorsWithoutRoot", func(t *testing.T) {
		t.Parallel()
		t.Run("provides all ancestor branches, oldest first", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.NewLineage()
			lineage.AddParent(three, two)
			lineage.AddParent(two, one)
			lineage.AddParent(one, main)
			have := lineage.AncestorsWithoutRoot(three)
			want := gitdomain.LocalBranchNames{one, two}
			must.Eq(t, want, have)
		})
		t.Run("one ancestor", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.NewLineage()
			lineage.AddParent(one, main)
			have := lineage.AncestorsWithoutRoot(one)
			want := gitdomain.LocalBranchNames{}
			must.Eq(t, want, have)
		})
		t.Run("no ancestors", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.NewLineage()
			lineage.AddParent(one, main)
			have := lineage.AncestorsWithoutRoot(two)
			want := gitdomain.LocalBranchNames{}
			must.Eq(t, want, have)
		})
	})

	t.Run("BranchAndAncestors", func(t *testing.T) {
		t.Parallel()
		lineage := configdomain.NewLineage()
		lineage.AddParent(one, main)
		have := lineage.BranchAndAncestors(one)
		want := gitdomain.LocalBranchNames{main, one}
		must.Eq(t, want, have)
	})

	t.Run("BranchLineageWithoutRoot", func(t *testing.T) {
		t.Parallel()
		t.Run("only root exists", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.Lineage{}
			have := lineage.BranchLineageWithoutRoot(main)
			want := gitdomain.LocalBranchNames{}
			must.Eq(t, want, have)
		})
		t.Run("one branch and root exist", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.NewLineage()
			lineage.AddParent(one, main)
			have := lineage.BranchLineageWithoutRoot(one)
			want := gitdomain.LocalBranchNames{one}
			must.Eq(t, want, have)
		})
		t.Run("multiple branches", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.NewLineage()
			lineage.AddParent(one, main)
			lineage.AddParent(two, one)
			want := gitdomain.LocalBranchNames{one, two}
			t.Run("given root", func(t *testing.T) {
				t.Parallel()
				have := lineage.BranchLineageWithoutRoot(main)
				must.Eq(t, want, have)
			})
			t.Run("given middle branch", func(t *testing.T) {
				t.Parallel()
				have := lineage.BranchLineageWithoutRoot(one)
				must.Eq(t, want, have)
			})
			t.Run("given leaf branch", func(t *testing.T) {
				t.Parallel()
				have := lineage.BranchLineageWithoutRoot(two)
				must.Eq(t, want, have)
			})
		})
	})

	t.Run("BranchNames", func(t *testing.T) {
		t.Parallel()
		lineage := configdomain.NewLineage()
		lineage.AddParent(one, main)
		lineage.AddParent(two, main)
		lineage.AddParent(three, main)
		have := lineage.BranchNames()
		want := gitdomain.LocalBranchNames{one, three, two}
		must.Eq(t, want, have)
	})

	t.Run("Branches", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.NewLineage()
			lineage.AddParent("branch-1", "branch-2")
			lineage.AddParent("branch-3", "branch-4")
			have := lineage.Branches()
			want := gitdomain.NewLocalBranchNames("branch-1", "branch-3")
			must.Eq(t, want, have)
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.NewLineage()
			have := lineage.Branches()
			want := gitdomain.LocalBranchNames{}
			must.Eq(t, want, have)
		})
	})

	t.Run("BranchesAndAncestors", func(t *testing.T) {
		t.Parallel()
		t.Run("deep lineage, multiple branches", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.NewLineage()
			lineage.AddParent(one, main)
			lineage.AddParent(two, one)
			lineage.AddParent(three, two)
			give := gitdomain.LocalBranchNames{two, one}
			have := lineage.BranchesAndAncestors(give)
			want := gitdomain.LocalBranchNames{main, one, two}
			must.Eq(t, want, have)
		})
		t.Run("deep lineage, multiple branches out of order", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.NewLineage()
			lineage.AddParent(one, main)
			lineage.AddParent(two, one)
			lineage.AddParent(three, two)
			give := gitdomain.LocalBranchNames{one, two, main, three}
			have := lineage.BranchesAndAncestors(give)
			want := gitdomain.LocalBranchNames{main, one, two, three}
			must.Eq(t, want, have)
		})
		t.Run("deep lineage, single branch", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.NewLineage()
			lineage.AddParent(one, main)
			lineage.AddParent(two, one)
			lineage.AddParent(three, two)
			give := gitdomain.LocalBranchNames{three}
			have := lineage.BranchesAndAncestors(give)
			want := gitdomain.LocalBranchNames{main, one, two, three}
			must.Eq(t, want, have)
		})
		t.Run("multiple lineages", func(t *testing.T) {
			t.Parallel()
			first := gitdomain.NewLocalBranchName("first")
			second := gitdomain.NewLocalBranchName("second")
			third := gitdomain.NewLocalBranchName("third")
			lineage := configdomain.NewLineage()
			lineage.AddParent(one, main)
			lineage.AddParent(two, one)
			lineage.AddParent(three, two)
			lineage.AddParent(first, main)
			lineage.AddParent(second, first)
			lineage.AddParent(third, second)
			give := gitdomain.LocalBranchNames{main, first, one, second, third, three, two}
			have := lineage.BranchesAndAncestors(give)
			want := gitdomain.LocalBranchNames{main, first, second, third, one, two, three}
			must.Eq(t, want, have)
		})
	})

	t.Run("Children", func(t *testing.T) {
		t.Parallel()
		t.Run("provides all children of the given branch, ordered alphabetically", func(t *testing.T) {
			t.Parallel()
			twoA := gitdomain.NewLocalBranchName("twoA")
			twoB := gitdomain.NewLocalBranchName("twoB")
			lineage := configdomain.NewLineage()
			lineage.AddParent(twoA, one)
			lineage.AddParent(twoB, one)
			have := lineage.Children(one)
			want := gitdomain.LocalBranchNames{twoA, twoB}
			must.Eq(t, want, have)
		})
		t.Run("provides only the immediate children, i.e. no grandchildren", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.NewLineage()
			lineage.AddParent(two, one)
			lineage.AddParent(three, two)
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
			lineage := configdomain.NewLineage()
			lineage.AddParent(two, one)
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
			lineage := configdomain.NewLineage()
			lineage.AddParent(branch, main)
			lineage.AddParent(other, main)
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
			lineage := configdomain.NewLineage()
			lineage.AddParent(branch, main)
			lineage.AddParent(child1, branch)
			lineage.AddParent(child2, branch)
			lineage.AddParent(other, main)
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
			lineage := configdomain.NewLineage()
			lineage.AddParent(branch, main)
			lineage.AddParent(child1, branch)
			lineage.AddParent(child1a, child1)
			lineage.AddParent(child1b, child1)
			lineage.AddParent(child2, branch)
			lineage.AddParent(child2a, child2)
			lineage.AddParent(child2b, child2)
			lineage.AddParent(other, main)
			have := lineage.Descendants(branch)
			want := gitdomain.LocalBranchNames{child1, child1a, child1b, child2, child2a, child2b}
			must.Eq(t, want, have)
		})
	})

	t.Run("IsAncestor", func(t *testing.T) {
		t.Run("recognizes greatgrandparent", func(t *testing.T) {
			t.Parallel()
			four := gitdomain.NewLocalBranchName("four")
			lineage := configdomain.NewLineage()
			lineage.AddParent(four, three)
			lineage.AddParent(three, two)
			lineage.AddParent(two, one)
			must.True(t, lineage.IsAncestor(one, four))
		})
		t.Run("child branches are not ancestors", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.NewLineage()
			lineage.AddParent(two, one)
			must.True(t, lineage.IsAncestor(one, two))
		})
		t.Run("unrelated branches are not ancestors", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.NewLineage()
			lineage.AddParent(two, one)
			lineage.AddParent(three, one)
			must.False(t, lineage.IsAncestor(two, three))
		})
	})

	t.Run("OrderHierarchically", func(t *testing.T) {
		t.Run("multiple lineages", func(t *testing.T) {
			t.Parallel()
			first := gitdomain.NewLocalBranchName("first")
			second := gitdomain.NewLocalBranchName("second")
			third := gitdomain.NewLocalBranchName("third")
			lineage := configdomain.NewLineage()
			lineage.AddParent(one, main)
			lineage.AddParent(two, one)
			lineage.AddParent(three, two)
			lineage.AddParent(first, main)
			lineage.AddParent(second, first)
			lineage.AddParent(third, second)
			give := lineage.BranchNames()
			want := gitdomain.LocalBranchNames{first, second, third, one, two, three}
			have := lineage.OrderHierarchically(give)
			must.Eq(t, want, have)
		})
		t.Run("single lineage", func(t *testing.T) {
			t.Parallel()
			four := gitdomain.NewLocalBranchName("four")
			lineage := configdomain.NewLineage()
			lineage.AddParent(one, main)
			lineage.AddParent(two, one)
			lineage.AddParent(three, two)
			lineage.AddParent(four, three)
			give := gitdomain.LocalBranchNames{four, one}
			want := gitdomain.LocalBranchNames{one, four}
			have := lineage.OrderHierarchically(give)
			must.Eq(t, want, have)
		})
		t.Run("elements out of order", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.NewLineage()
			lineage.AddParent(one, main)
			lineage.AddParent(two, one)
			lineage.AddParent(three, two)
			give := gitdomain.LocalBranchNames{one, two, main, three}
			want := gitdomain.LocalBranchNames{main, one, two, three}
			have := lineage.OrderHierarchically(give)
			must.Eq(t, want, have)
		})
		t.Run("perennial branches", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.Lineage{}
			give := gitdomain.LocalBranchNames{one, two, main}
			want := gitdomain.LocalBranchNames{one, two, main}
			have := lineage.OrderHierarchically(give)
			must.Eq(t, want, have)
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.Lineage{}
			give := gitdomain.LocalBranchNames{}
			have := lineage.OrderHierarchically(give)
			must.Eq(t, 0, len(have))
		})
	})

	t.Run("Parent", func(t *testing.T) {
		t.Parallel()
		t.Run("feature branch", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.NewLineage()
			lineage.AddParent(one, main)
			have := lineage.Parent(one)
			want := Some(main)
			must.Eq(t, want, have)
		})
		t.Run("main branch", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.Lineage{}
			have := lineage.Parent(main)
			want := None[gitdomain.LocalBranchName]()
			must.EqOp(t, want, have)
		})
		t.Run("perennial branch", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.Lineage{}
			have := lineage.Parent(one)
			want := None[gitdomain.LocalBranchName]()
			must.EqOp(t, want, have)
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
			have := configdomain.NewLineage()
			have.AddParent(branch1, main)
			have.AddParent(branch1a, branch1)
			have.AddParent(branch1b, branch1)
			have.AddParent(branch2, main)
			have.RemoveBranch(branch1)
			want := configdomain.NewLineage()
			want.AddParent(branch1a, main)
			want.AddParent(branch1b, main)
			want.AddParent(branch2, main)
			must.Eq(t, want, have)
		})
		t.Run("branch is a child branch", func(t *testing.T) {
			t.Parallel()
			main := gitdomain.NewLocalBranchName("main")
			branch1 := gitdomain.NewLocalBranchName("branch-1")
			branch2 := gitdomain.NewLocalBranchName("branch-2")
			have := configdomain.NewLineage()
			have.AddParent(branch1, main)
			have.AddParent(branch2, main)
			have.RemoveBranch(branch1)
			want := configdomain.NewLineage()
			want.AddParent(branch2, main)
			must.Eq(t, want, have)
		})
		t.Run("branch is not in lineage", func(t *testing.T) {
			t.Parallel()
			main := gitdomain.NewLocalBranchName("main")
			branch1 := gitdomain.NewLocalBranchName("branch-1")
			branch2 := gitdomain.NewLocalBranchName("branch-2")
			have := configdomain.NewLineage()
			have.AddParent(branch1, main)
			have.RemoveBranch(branch2)
			want := configdomain.NewLineage()
			want.AddParent(branch1, main)
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
			lineage := configdomain.NewLineage()
			lineage.AddParent(two, one)
			lineage.AddParent(one, main)
			lineage.AddParent(hotfix1, prod)
			lineage.AddParent(hotfix2, prod)
			have := lineage.Roots()
			want := gitdomain.LocalBranchNames{main, prod}
			must.Eq(t, want, have)
		})
		t.Run("no stacked changes", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.NewLineage()
			lineage.AddParent(one, main)
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
