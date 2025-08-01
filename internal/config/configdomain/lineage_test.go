package configdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
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
			lineage := configdomain.NewLineageWith(configdomain.LineageData{
				three: two,
				two:   one,
				one:   main,
			})
			have := lineage.Ancestors(three)
			want := gitdomain.LocalBranchNames{main, one, two}
			must.Eq(t, want, have)
		})
		t.Run("one ancestor", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.NewLineageWith(configdomain.LineageData{
				one: main,
			})
			have := lineage.Ancestors(one)
			want := gitdomain.LocalBranchNames{main}
			must.Eq(t, want, have)
		})
		t.Run("no ancestors", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.NewLineageWith(configdomain.LineageData{
				one: main,
			})
			have := lineage.Ancestors(two)
			want := gitdomain.LocalBranchNames{}
			must.Eq(t, want, have)
		})
	})

	t.Run("AncestorsWithoutRoot", func(t *testing.T) {
		t.Parallel()
		t.Run("provides all ancestor branches, oldest first", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.NewLineageWith(configdomain.LineageData{
				three: two,
				two:   one,
				one:   main,
			})
			have := lineage.AncestorsWithoutRoot(three)
			want := gitdomain.LocalBranchNames{one, two}
			must.Eq(t, want, have)
		})
		t.Run("one ancestor", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.NewLineageWith(configdomain.LineageData{
				one: main,
			})
			have := lineage.AncestorsWithoutRoot(one)
			want := gitdomain.LocalBranchNames{}
			must.Eq(t, want, have)
		})
		t.Run("no ancestors", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.NewLineageWith(configdomain.LineageData{
				one: main,
			})
			have := lineage.AncestorsWithoutRoot(two)
			want := gitdomain.LocalBranchNames{}
			must.Eq(t, want, have)
		})
	})

	t.Run("BranchAndAncestors", func(t *testing.T) {
		t.Parallel()
		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			one: main,
		})
		have := lineage.BranchAndAncestors(one)
		want := gitdomain.LocalBranchNames{main, one}
		must.Eq(t, want, have)
	})

	t.Run("BranchAndAncestorsWithoutRoot", func(t *testing.T) {
		t.Parallel()
		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			one: main,
			two: one,
		})
		have := lineage.BranchAndAncestorsWithoutRoot(two)
		want := gitdomain.LocalBranchNames{one, two}
		must.Eq(t, want, have)
	})

	t.Run("BranchLineageWithoutRoot", func(t *testing.T) {
		t.Parallel()
		t.Run("only root exists", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.Lineage{}
			have := lineage.BranchLineageWithoutRoot(main, gitdomain.LocalBranchNames{main})
			want := gitdomain.LocalBranchNames{}
			must.Eq(t, want, have)
		})
		t.Run("one branch and root exist", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.NewLineageWith(configdomain.LineageData{
				one: main,
			})
			have := lineage.BranchLineageWithoutRoot(one, gitdomain.LocalBranchNames{main})
			want := gitdomain.LocalBranchNames{one}
			must.Eq(t, want, have)
		})
		t.Run("multiple branches", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.NewLineageWith(configdomain.LineageData{
				one: main,
				two: one,
			})
			want := gitdomain.LocalBranchNames{one, two}
			t.Run("given root", func(t *testing.T) {
				t.Parallel()
				have := lineage.BranchLineageWithoutRoot(main, gitdomain.LocalBranchNames{main})
				must.Eq(t, want, have)
			})
			t.Run("given middle branch", func(t *testing.T) {
				t.Parallel()
				have := lineage.BranchLineageWithoutRoot(one, gitdomain.LocalBranchNames{main})
				must.Eq(t, want, have)
			})
			t.Run("given leaf branch", func(t *testing.T) {
				t.Parallel()
				have := lineage.BranchLineageWithoutRoot(two, gitdomain.LocalBranchNames{main})
				must.Eq(t, want, have)
			})
		})
		t.Run("branch without an ancestor", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.Lineage{}
			have := lineage.BranchLineageWithoutRoot(one, gitdomain.LocalBranchNames{main})
			want := gitdomain.LocalBranchNames{one}
			must.Eq(t, want, have)
		})
	})

	t.Run("BranchNames", func(t *testing.T) {
		t.Parallel()
		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			one:   main,
			two:   main,
			three: main,
		})
		have := lineage.BranchNames()
		want := gitdomain.LocalBranchNames{one, three, two}
		must.Eq(t, want, have)
	})

	t.Run("BranchesAndAncestors", func(t *testing.T) {
		t.Parallel()
		t.Run("deep lineage, multiple branches", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.NewLineageWith(configdomain.LineageData{
				one:   main,
				two:   one,
				three: two,
			})
			give := gitdomain.LocalBranchNames{two, one}
			have := lineage.BranchesAndAncestors(give)
			want := gitdomain.LocalBranchNames{main, one, two}
			must.Eq(t, want, have)
		})
		t.Run("deep lineage, multiple branches out of order", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.NewLineageWith(configdomain.LineageData{
				one:   main,
				two:   one,
				three: two,
			})
			give := gitdomain.LocalBranchNames{one, two, main, three}
			have := lineage.BranchesAndAncestors(give)
			want := gitdomain.LocalBranchNames{main, one, two, three}
			must.Eq(t, want, have)
		})
		t.Run("deep lineage, single branch", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.NewLineageWith(configdomain.LineageData{
				one:   main,
				two:   one,
				three: two,
			})
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
			lineage := configdomain.NewLineageWith(configdomain.LineageData{
				one:    main,
				two:    one,
				three:  two,
				first:  main,
				second: first,
				third:  second,
			})
			give := gitdomain.LocalBranchNames{main, first, one, second, third, three, two}
			have := lineage.BranchesAndAncestors(give)
			want := gitdomain.LocalBranchNames{main, first, second, third, one, two, three}
			must.Eq(t, want, have)
		})
	})

	t.Run("BranchesWithParents", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.NewLineageWith(configdomain.LineageData{
				"branch-1": "branch-2",
				"branch-3": "branch-4",
			})
			have := lineage.BranchesWithParents()
			want := gitdomain.NewLocalBranchNames("branch-1", "branch-3")
			must.Eq(t, want, have)
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.NewLineage()
			have := lineage.BranchesWithParents()
			want := gitdomain.LocalBranchNames(nil)
			must.Eq(t, want, have)
		})
	})

	t.Run("Children", func(t *testing.T) {
		t.Parallel()
		t.Run("provides all children of the given branch, ordered alphabetically", func(t *testing.T) {
			t.Parallel()
			twoA := gitdomain.NewLocalBranchName("twoA")
			twoB := gitdomain.NewLocalBranchName("twoB")
			lineage := configdomain.NewLineageWith(configdomain.LineageData{
				twoA: one,
				twoB: one,
			})
			have := lineage.Children(one)
			want := gitdomain.LocalBranchNames{twoA, twoB}
			must.Eq(t, want, have)
		})
		t.Run("provides only the immediate children, i.e. no grandchildren", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.NewLineageWith(configdomain.LineageData{
				two:   one,
				three: two,
			})
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

	t.Run("Descendants", func(t *testing.T) {
		t.Parallel()
		t.Run("branch has no children", func(t *testing.T) {
			t.Parallel()
			branch := gitdomain.NewLocalBranchName("branch")
			other := gitdomain.NewLocalBranchName("other")
			lineage := configdomain.NewLineageWith(configdomain.LineageData{
				branch: main,
				other:  main,
			})
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
			lineage := configdomain.NewLineageWith(configdomain.LineageData{
				branch: main,
				child1: branch,
				child2: branch,
				other:  main,
			})
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
			lineage := configdomain.NewLineageWith(configdomain.LineageData{
				branch:  main,
				child1:  branch,
				child1a: child1,
				child1b: child1,
				child2:  branch,
				child2a: child2,
				child2b: child2,
				other:   main,
			})
			have := lineage.Descendants(branch)
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
			lineage := configdomain.NewLineageWith(configdomain.LineageData{
				branch1: branch2,
				branch3: branch4,
			})
			have := lineage.Entries()
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

	t.Run("HasDescendents", func(t *testing.T) {
		t.Parallel()
		t.Run("has a descendent", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.NewLineageWith(configdomain.LineageData{
				one: two,
			})
			must.True(t, lineage.HasDescendents(two))
		})
		t.Run("has no descendent", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.NewLineageWith(configdomain.LineageData{
				one: two,
			})
			must.False(t, lineage.HasDescendents(one))
		})
	})

	t.Run("HasParents", func(t *testing.T) {
		t.Parallel()
		t.Run("has a parent", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.NewLineageWith(configdomain.LineageData{
				two: one,
			})
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
			four := gitdomain.NewLocalBranchName("four")
			lineage := configdomain.NewLineageWith(configdomain.LineageData{
				four:  three,
				three: two,
				two:   one,
			})
			must.True(t, lineage.IsAncestor(one, four))
		})
		t.Run("child branches are not ancestors", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.NewLineageWith(configdomain.LineageData{
				two: one,
			})
			must.True(t, lineage.IsAncestor(one, two))
		})
		t.Run("unrelated branches are not ancestors", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.NewLineageWith(configdomain.LineageData{
				two:   one,
				three: one,
			})
			must.False(t, lineage.IsAncestor(two, three))
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
			lineage := configdomain.NewLineageWith(configdomain.LineageData{
				"branch-1": "branch-2",
			})
			must.False(t, lineage.IsEmpty())
		})
	})

	t.Run("LatestAncestor", func(t *testing.T) {
		t.Parallel()
		t.Run("happy path", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.NewLineageWith(configdomain.LineageData{
				one:   main,
				two:   one,
				three: two,
			})
			have := lineage.LatestAncestor(three, gitdomain.LocalBranchNames{one, two})
			must.Eq(t, Some(two), have)
		})
		t.Run("no candidates", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.NewLineageWith(configdomain.LineageData{
				one:   main,
				two:   one,
				three: two,
			})
			have := lineage.LatestAncestor(three, gitdomain.LocalBranchNames{})
			must.Eq(t, None[gitdomain.LocalBranchName](), have)
		})
		t.Run("candidates contains branch", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.NewLineageWith(configdomain.LineageData{
				one:   main,
				two:   one,
				three: two,
			})
			have := lineage.LatestAncestor(three, gitdomain.LocalBranchNames{two, three})
			must.Eq(t, Some(three), have)
		})
		t.Run("candidates not in lineage", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.NewLineageWith(configdomain.LineageData{
				one: main,
			})
			have := lineage.LatestAncestor(two, gitdomain.LocalBranchNames{three})
			must.Eq(t, None[gitdomain.LocalBranchName](), have)
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
			lineage := configdomain.NewLineageWith(configdomain.LineageData{
				"branch-1": "branch-2",
				"branch-3": "branch-4",
			})
			must.EqOp(t, 2, lineage.Len())
		})
	})

	t.Run("Merge", func(t *testing.T) {
		t.Parallel()
		lineage1 := configdomain.NewLineageWith(configdomain.LineageData{
			"one": "main",
			"two": "one",
		})
		lineage2 := configdomain.NewLineageWith(configdomain.LineageData{
			"alpha": "main",
			"beta":  "alpha",
		})
		haveMerged := lineage1.Merge(lineage2)
		want := configdomain.NewLineageWith(configdomain.LineageData{
			"one":   "main",
			"two":   "one",
			"alpha": "main",
			"beta":  "alpha",
		})
		must.Eq(t, want, haveMerged)
	})

	t.Run("OrderHierarchically", func(t *testing.T) {
		t.Run("multiple lineages", func(t *testing.T) {
			t.Parallel()
			first := gitdomain.NewLocalBranchName("first")
			second := gitdomain.NewLocalBranchName("second")
			third := gitdomain.NewLocalBranchName("third")
			lineage := configdomain.NewLineageWith(configdomain.LineageData{
				one:    main,
				two:    one,
				three:  two,
				first:  main,
				second: first,
				third:  second,
			})
			give := lineage.BranchNames()
			want := gitdomain.LocalBranchNames{first, second, third, one, two, three}
			have := lineage.OrderHierarchically(give)
			must.Eq(t, want, have)
		})
		t.Run("single lineage", func(t *testing.T) {
			t.Parallel()
			four := gitdomain.NewLocalBranchName("four")
			lineage := configdomain.NewLineageWith(configdomain.LineageData{
				one:   main,
				two:   one,
				three: two,
				four:  three,
			})
			give := gitdomain.LocalBranchNames{four, one}
			want := gitdomain.LocalBranchNames{one, four}
			have := lineage.OrderHierarchically(give)
			must.Eq(t, want, have)
		})
		t.Run("elements out of order", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.NewLineageWith(configdomain.LineageData{
				one:   main,
				two:   one,
				three: two,
			})
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
			lineage := configdomain.NewLineageWith(configdomain.LineageData{
				one: main,
			})
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
			have := configdomain.NewLineageWith(configdomain.LineageData{
				branch1:  main,
				branch1a: branch1,
				branch1b: branch1,
				branch2:  main,
			})
			have = have.RemoveBranch(branch1)
			want := configdomain.NewLineageWith(configdomain.LineageData{
				branch1a: branch1,
				branch1b: branch1,
				branch2:  main,
			})
			must.Eq(t, want, have)
		})
		t.Run("branch is a child branch", func(t *testing.T) {
			t.Parallel()
			main := gitdomain.NewLocalBranchName("main")
			branch1 := gitdomain.NewLocalBranchName("branch-1")
			branch2 := gitdomain.NewLocalBranchName("branch-2")
			have := configdomain.NewLineageWith(configdomain.LineageData{
				branch1: main,
				branch2: main,
			})
			have = have.RemoveBranch(branch1)
			want := configdomain.NewLineageWith(configdomain.LineageData{
				branch2: main,
			})
			must.Eq(t, want, have)
		})
		t.Run("branch is not in lineage", func(t *testing.T) {
			t.Parallel()
			main := gitdomain.NewLocalBranchName("main")
			branch1 := gitdomain.NewLocalBranchName("branch-1")
			branch2 := gitdomain.NewLocalBranchName("branch-2")
			have := configdomain.NewLineageWith(configdomain.LineageData{
				branch1: main,
			})
			have = have.RemoveBranch(branch2)
			want := configdomain.NewLineageWith(configdomain.LineageData{
				branch1: main,
			})
			must.Eq(t, want, have)
		})
	})

	t.Run("Root", func(t *testing.T) {
		t.Parallel()
		t.Run("stack", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.NewLineageWith(configdomain.LineageData{
				three: two,
				two:   one,
				one:   main,
			})
			have := lineage.Root(three)
			must.Eq(t, main, have)
		})
		t.Run("one ancestor", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.NewLineageWith(configdomain.LineageData{
				one: main,
			})
			have := lineage.Root(one)
			must.Eq(t, main, have)
		})
		t.Run("no ancestors", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.NewLineageWith(configdomain.LineageData{})
			have := lineage.Root(two)
			must.Eq(t, two, have)
		})
	})

	t.Run("Roots", func(t *testing.T) {
		t.Parallel()
		t.Run("multiple roots with child branches", func(t *testing.T) {
			t.Parallel()
			prod := gitdomain.NewLocalBranchName("prod")
			hotfix1 := gitdomain.NewLocalBranchName("hotfix1")
			hotfix2 := gitdomain.NewLocalBranchName("hotfix2")
			lineage := configdomain.NewLineageWith(configdomain.LineageData{
				two:     one,
				one:     main,
				hotfix1: prod,
				hotfix2: prod,
			})
			have := lineage.Roots()
			want := gitdomain.LocalBranchNames{main, prod}
			must.Eq(t, want, have)
		})
		t.Run("no stacked changes", func(t *testing.T) {
			t.Parallel()
			lineage := configdomain.NewLineageWith(configdomain.LineageData{
				one: main,
			})
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

	t.Run("Set", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			branch := gitdomain.NewLocalBranchName("branch")
			parent := gitdomain.NewLocalBranchName("parent")
			lineage := configdomain.NewLineage()
			lineage = lineage.Set(branch, parent)
			have, has := lineage.Parent(branch).Get()
			must.True(t, has)
			must.Eq(t, parent, have)
		})
		t.Run("entry already exists", func(t *testing.T) {
			t.Parallel()
			branch := gitdomain.NewLocalBranchName("branch")
			parent := gitdomain.NewLocalBranchName("parent")
			lineage := configdomain.NewLineageWith(configdomain.LineageData{
				branch: parent,
			})
			lineage = lineage.Set(branch, parent)
			must.EqOp(t, 1, lineage.Len())
			have, has := lineage.Parent(branch).Get()
			must.True(t, has)
			must.Eq(t, parent, have)
		})
	})
}
