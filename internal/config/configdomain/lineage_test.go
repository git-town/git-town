package configdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	. "github.com/git-town/git-town/v16/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestLineage(t *testing.T) {
	t.Parallel()
	main := gitdomain.NewLocalBranchName("main")
	one := gitdomain.NewLocalBranchName("one")
	two := gitdomain.NewLocalBranchName("two")
	three := gitdomain.NewLocalBranchName("three")

	t.Run("Add", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			builder := configdomain.NewLineageBuilder()
			branch := gitdomain.NewLocalBranchName("branch")
			parent := gitdomain.NewLocalBranchName("parent")
			builder.Add(branch, parent)
			have, has := builder.Lineage().Parent(branch).Get()
			must.True(t, has)
			must.Eq(t, parent, have)
		})
		t.Run("entry already exists", func(t *testing.T) {
			t.Parallel()
			builder := configdomain.NewLineageBuilder()
			branch := gitdomain.NewLocalBranchName("branch")
			parent := gitdomain.NewLocalBranchName("parent")
			builder.Add(branch, parent)
			builder.Add(branch, parent)
			lineage := builder.Lineage()
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
			builder := configdomain.NewLineageBuilder()
			builder.Add(three, two)
			builder.Add(two, one)
			builder.Add(one, main)
			have := builder.Lineage().Ancestors(three)
			want := gitdomain.LocalBranchNames{main, one, two}
			must.Eq(t, want, have)
		})
		t.Run("one ancestor", func(t *testing.T) {
			t.Parallel()
			builder := configdomain.NewLineageBuilder()
			builder.Add(one, main)
			have := builder.Lineage().Ancestors(one)
			want := gitdomain.LocalBranchNames{main}
			must.Eq(t, want, have)
		})
		t.Run("no ancestors", func(t *testing.T) {
			t.Parallel()
			builder := configdomain.NewLineageBuilder()
			builder.Add(one, main)
			have := builder.Lineage().Ancestors(two)
			want := gitdomain.LocalBranchNames{}
			must.Eq(t, want, have)
		})
	})

	t.Run("AncestorsWithoutRoot", func(t *testing.T) {
		t.Parallel()
		t.Run("provides all ancestor branches, oldest first", func(t *testing.T) {
			t.Parallel()
			builder := configdomain.NewLineageBuilder()
			builder.Add(three, two)
			builder.Add(two, one)
			builder.Add(one, main)
			have := builder.Lineage().AncestorsWithoutRoot(three)
			want := gitdomain.LocalBranchNames{one, two}
			must.Eq(t, want, have)
		})
		t.Run("one ancestor", func(t *testing.T) {
			t.Parallel()
			builder := configdomain.NewLineageBuilder()
			builder.Add(one, main)
			have := builder.Lineage().AncestorsWithoutRoot(one)
			want := gitdomain.LocalBranchNames{}
			must.Eq(t, want, have)
		})
		t.Run("no ancestors", func(t *testing.T) {
			t.Parallel()
			builder := configdomain.NewLineageBuilder()
			builder.Add(one, main)
			have := builder.Lineage().AncestorsWithoutRoot(two)
			want := gitdomain.LocalBranchNames{}
			must.Eq(t, want, have)
		})
	})

	t.Run("BranchAndAncestors", func(t *testing.T) {
		t.Parallel()
		builder := configdomain.NewLineageBuilder()
		builder.Add(one, main)
		have := builder.Lineage().BranchAndAncestors(one)
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
			builder := configdomain.NewLineageBuilder()
			builder.Add(one, main)
			have := builder.Lineage().BranchLineageWithoutRoot(one)
			want := gitdomain.LocalBranchNames{one}
			must.Eq(t, want, have)
		})
		t.Run("multiple branches", func(t *testing.T) {
			t.Parallel()
			builder := configdomain.NewLineageBuilder()
			builder.Add(one, main)
			builder.Add(two, one)
			want := gitdomain.LocalBranchNames{one, two}
			lineage := builder.Lineage()
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
		builder := configdomain.NewLineageBuilder()
		builder.Add(one, main)
		builder.Add(two, main)
		builder.Add(three, main)
		have := builder.Lineage().BranchNames()
		want := gitdomain.LocalBranchNames{one, three, two}
		must.Eq(t, want, have)
	})

	t.Run("BranchesWithParents", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			builder := configdomain.NewLineageBuilder()
			builder.Add("branch-1", "branch-2")
			builder.Add("branch-3", "branch-4")
			have := builder.Lineage().BranchesWithParents()
			want := gitdomain.NewLocalBranchNames("branch-1", "branch-3")
			must.Eq(t, want, have)
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.NewLineage()
			have := lineage.BranchesWithParents()
			want := gitdomain.LocalBranchNames{}
			must.Eq(t, want, have)
		})
	})

	t.Run("BranchesAndAncestors", func(t *testing.T) {
		t.Parallel()
		t.Run("deep lineage, multiple branches", func(t *testing.T) {
			t.Parallel()
			builder := configdomain.NewLineageBuilder()
			builder.Add(one, main)
			builder.Add(two, one)
			builder.Add(three, two)
			give := gitdomain.LocalBranchNames{two, one}
			have := builder.Lineage().BranchesAndAncestors(give)
			want := gitdomain.LocalBranchNames{main, one, two}
			must.Eq(t, want, have)
		})
		t.Run("deep lineage, multiple branches out of order", func(t *testing.T) {
			t.Parallel()
			builder := configdomain.NewLineageBuilder()
			builder.Add(one, main)
			builder.Add(two, one)
			builder.Add(three, two)
			give := gitdomain.LocalBranchNames{one, two, main, three}
			have := builder.Lineage().BranchesAndAncestors(give)
			want := gitdomain.LocalBranchNames{main, one, two, three}
			must.Eq(t, want, have)
		})
		t.Run("deep lineage, single branch", func(t *testing.T) {
			t.Parallel()
			builder := configdomain.NewLineageBuilder()
			builder.Add(one, main)
			builder.Add(two, one)
			builder.Add(three, two)
			give := gitdomain.LocalBranchNames{three}
			have := builder.Lineage().BranchesAndAncestors(give)
			want := gitdomain.LocalBranchNames{main, one, two, three}
			must.Eq(t, want, have)
		})
		t.Run("multiple lineages", func(t *testing.T) {
			t.Parallel()
			first := gitdomain.NewLocalBranchName("first")
			second := gitdomain.NewLocalBranchName("second")
			third := gitdomain.NewLocalBranchName("third")
			builder := configdomain.NewLineageBuilder()
			builder.Add(one, main)
			builder.Add(two, one)
			builder.Add(three, two)
			builder.Add(first, main)
			builder.Add(second, first)
			builder.Add(third, second)
			give := gitdomain.LocalBranchNames{main, first, one, second, third, three, two}
			have := builder.Lineage().BranchesAndAncestors(give)
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
			builder := configdomain.NewLineageBuilder()
			builder.Add(twoA, one)
			builder.Add(twoB, one)
			have := builder.Lineage().Children(one)
			want := gitdomain.LocalBranchNames{twoA, twoB}
			must.Eq(t, want, have)
		})
		t.Run("provides only the immediate children, i.e. no grandchildren", func(t *testing.T) {
			t.Parallel()
			builder := configdomain.NewLineageBuilder()
			builder.Add(two, one)
			builder.Add(three, two)
			have := builder.Lineage().Children(one)
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

	t.Run("Descendants", func(t *testing.T) {
		t.Parallel()
		t.Run("branch has no children", func(t *testing.T) {
			t.Parallel()
			branch := gitdomain.NewLocalBranchName("branch")
			other := gitdomain.NewLocalBranchName("other")
			builder := configdomain.NewLineageBuilder()
			builder.Add(branch, main)
			builder.Add(other, main)
			have := builder.Lineage().Descendants(branch)
			want := gitdomain.LocalBranchNames{}
			must.Eq(t, want, have)
		})
		t.Run("branch has only direct children", func(t *testing.T) {
			t.Parallel()
			branch := gitdomain.NewLocalBranchName("branch")
			child1 := gitdomain.NewLocalBranchName("child1")
			child2 := gitdomain.NewLocalBranchName("child2")
			other := gitdomain.NewLocalBranchName("other")
			builder := configdomain.NewLineageBuilder()
			builder.Add(branch, main)
			builder.Add(child1, branch)
			builder.Add(child2, branch)
			builder.Add(other, main)
			have := builder.Lineage().Descendants(branch)
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
			builder := configdomain.NewLineageBuilder()
			builder.Add(branch, main)
			builder.Add(child1, branch)
			builder.Add(child1a, child1)
			builder.Add(child1b, child1)
			builder.Add(child2, branch)
			builder.Add(child2a, child2)
			builder.Add(child2b, child2)
			builder.Add(other, main)
			have := builder.Lineage().Descendants(branch)
			want := gitdomain.LocalBranchNames{child1, child1a, child1b, child2, child2a, child2b}
			must.Eq(t, want, have)
		})
	})

	t.Run("Entries", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			branch1 := gitdomain.NewLocalBranchName("branch-1")
			branch2 := gitdomain.NewLocalBranchName("branch-2")
			branch3 := gitdomain.NewLocalBranchName("branch-3")
			branch4 := gitdomain.NewLocalBranchName("branch-4")
			builder := configdomain.NewLineageBuilder()
			builder.Add(branch1, branch2)
			builder.Add(branch3, branch4)
			have := builder.Lineage().Entries()
			want := []configdomain.LineageEntry{
				{
					Child:  branch1,
					Parent: branch2,
				},
				{
					Child:  branch3,
					Parent: branch4,
				},
			}
			must.Eq(t, want, have)
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.NewLineage()
			have := lineage.Entries()
			want := []configdomain.LineageEntry{}
			must.Eq(t, want, have)
		})
	})

	t.Run("HasParents", func(t *testing.T) {
		t.Parallel()
		t.Run("has a parent", func(t *testing.T) {
			t.Parallel()
			builder := configdomain.NewLineageBuilder()
			builder.Add(two, one)
			must.True(t, builder.Lineage().HasParents(two))
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
			four := gitdomain.NewLocalBranchName("four")
			builder := configdomain.NewLineageBuilder()
			builder.Add(four, three)
			builder.Add(three, two)
			builder.Add(two, one)
			must.True(t, builder.Lineage().IsAncestor(one, four))
		})
		t.Run("child branches are not ancestors", func(t *testing.T) {
			t.Parallel()
			builder := configdomain.NewLineageBuilder()
			builder.Add(two, one)
			must.True(t, builder.Lineage().IsAncestor(one, two))
		})
		t.Run("unrelated branches are not ancestors", func(t *testing.T) {
			t.Parallel()
			builder := configdomain.NewLineageBuilder()
			builder.Add(two, one)
			builder.Add(three, one)
			must.False(t, builder.Lineage().IsAncestor(two, three))
		})
	})

	t.Run("IsEmpty", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.NewLineage()
			must.True(t, lineage.IsEmpty())
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			builder := configdomain.NewLineageBuilder()
			builder.Add("branch-1", "branch-2")
			must.False(t, builder.Lineage().IsEmpty())
		})
	})

	t.Run("Len", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.NewLineage()
			must.EqOp(t, 0, lineage.Len())
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			builder := configdomain.NewLineageBuilder()
			builder.Add("branch-1", "branch-2")
			builder.Add("branch-3", "branch-4")
			must.EqOp(t, 2, builder.Lineage().Len())
		})
	})

	t.Run("Merge", func(t *testing.T) {
		t.Parallel()
		builder1 := configdomain.NewLineageBuilder()
		builder1.Add("one", "main")
		builder1.Add("two", "one")
		builder2 := configdomain.NewLineageBuilder()
		builder2.Add("alpha", "main")
		builder2.Add("beta", "alpha")
		haveMerged := builder1.Lineage().Merge(builder2.Lineage())
		wantBuilder := configdomain.NewLineageBuilder()
		wantBuilder.Add("one", "main")
		wantBuilder.Add("two", "one")
		wantBuilder.Add("alpha", "main")
		wantBuilder.Add("beta", "alpha")
		must.Eq(t, wantBuilder.Lineage(), haveMerged)
	})

	t.Run("OrderHierarchically", func(t *testing.T) {
		t.Run("multiple lineages", func(t *testing.T) {
			t.Parallel()
			first := gitdomain.NewLocalBranchName("first")
			second := gitdomain.NewLocalBranchName("second")
			third := gitdomain.NewLocalBranchName("third")
			builder := configdomain.NewLineageBuilder()
			builder.Add(one, main)
			builder.Add(two, one)
			builder.Add(three, two)
			builder.Add(first, main)
			builder.Add(second, first)
			builder.Add(third, second)
			lineage := builder.Lineage()
			give := lineage.BranchNames()
			want := gitdomain.LocalBranchNames{first, second, third, one, two, three}
			have := lineage.OrderHierarchically(give)
			must.Eq(t, want, have)
		})
		t.Run("single lineage", func(t *testing.T) {
			t.Parallel()
			four := gitdomain.NewLocalBranchName("four")
			builder := configdomain.NewLineageBuilder()
			builder.Add(one, main)
			builder.Add(two, one)
			builder.Add(three, two)
			builder.Add(four, three)
			give := gitdomain.LocalBranchNames{four, one}
			want := gitdomain.LocalBranchNames{one, four}
			have := builder.Lineage().OrderHierarchically(give)
			must.Eq(t, want, have)
		})
		t.Run("elements out of order", func(t *testing.T) {
			t.Parallel()
			builder := configdomain.NewLineageBuilder()
			builder.Add(one, main)
			builder.Add(two, one)
			builder.Add(three, two)
			give := gitdomain.LocalBranchNames{one, two, main, three}
			want := gitdomain.LocalBranchNames{main, one, two, three}
			have := builder.Lineage().OrderHierarchically(give)
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
			builder := configdomain.NewLineageBuilder()
			builder.Add(one, main)
			have := builder.Lineage().Parent(one)
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
			haveBuilder := configdomain.NewLineageBuilder()
			haveBuilder.Add(branch1, main)
			haveBuilder.Add(branch1a, branch1)
			haveBuilder.Add(branch1b, branch1)
			haveBuilder.Add(branch2, main)
			have := haveBuilder.Lineage()
			have = have.RemoveBranch(branch1)
			wantBuilder := configdomain.NewLineageBuilder()
			wantBuilder.Add(branch1a, branch1)
			wantBuilder.Add(branch1b, branch1)
			wantBuilder.Add(branch2, main)
			must.Eq(t, wantBuilder.Lineage(), have)
		})
		t.Run("branch is a child branch", func(t *testing.T) {
			t.Parallel()
			main := gitdomain.NewLocalBranchName("main")
			branch1 := gitdomain.NewLocalBranchName("branch-1")
			branch2 := gitdomain.NewLocalBranchName("branch-2")
			haveBuilder := configdomain.NewLineageBuilder()
			haveBuilder.Add(branch1, main)
			haveBuilder.Add(branch2, main)
			have := haveBuilder.Lineage()
			have = have.RemoveBranch(branch1)
			wantBuilder := configdomain.NewLineageBuilder()
			wantBuilder.Add(branch2, main)
			must.Eq(t, wantBuilder.Lineage(), have)
		})
		t.Run("branch is not in lineage", func(t *testing.T) {
			t.Parallel()
			main := gitdomain.NewLocalBranchName("main")
			branch1 := gitdomain.NewLocalBranchName("branch-1")
			branch2 := gitdomain.NewLocalBranchName("branch-2")
			haveBuilder := configdomain.NewLineageBuilder()
			haveBuilder.Add(branch1, main)
			have := haveBuilder.Lineage()
			have = have.RemoveBranch(branch2)
			wantBuilder := configdomain.NewLineageBuilder()
			wantBuilder.Add(branch1, main)
			must.Eq(t, wantBuilder.Lineage(), have)
		})
	})

	t.Run("Roots", func(t *testing.T) {
		t.Parallel()
		t.Run("multiple roots with child branches", func(t *testing.T) {
			t.Parallel()
			prod := gitdomain.NewLocalBranchName("prod")
			hotfix1 := gitdomain.NewLocalBranchName("hotfix1")
			hotfix2 := gitdomain.NewLocalBranchName("hotfix2")
			builder := configdomain.NewLineageBuilder()
			builder.Add(two, one)
			builder.Add(one, main)
			builder.Add(hotfix1, prod)
			builder.Add(hotfix2, prod)
			have := builder.Lineage().Roots()
			want := gitdomain.LocalBranchNames{main, prod}
			must.Eq(t, want, have)
		})
		t.Run("no stacked changes", func(t *testing.T) {
			t.Parallel()
			builder := configdomain.NewLineageBuilder()
			builder.Add(one, main)
			have := builder.Lineage().Roots()
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
