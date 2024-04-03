package gitdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v13/src/git/gitdomain"
	"github.com/shoenig/test/must"
)

func TestLocalBranchNames(t *testing.T) {
	t.Parallel()

	t.Run("AppendAllMissing", func(t *testing.T) {
		t.Parallel()
		t.Run("append some new elements", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.NewLocalBranchNames("one", "two")
			have := give.AppendAllMissing("two", "three", "four")
			want := gitdomain.NewLocalBranchNames("one", "two", "three", "four")
			must.Eq(t, want, have)
		})
		t.Run("append no new elements", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.NewLocalBranchNames("one", "two")
			have := give.AppendAllMissing("one", "two")
			want := gitdomain.NewLocalBranchNames("one", "two")
			must.Eq(t, want, have)
		})
	})

	t.Run("AtRemote", func(t *testing.T) {
		t.Parallel()
		branch := gitdomain.NewLocalBranchName("branch")
		have := branch.AtRemote(gitdomain.RemoteOrigin)
		want := gitdomain.NewRemoteBranchName("origin/branch")
		must.EqOp(t, want, have)
	})

	t.Run("Contains", func(t *testing.T) {
		t.Parallel()
		branches := gitdomain.NewLocalBranchNames("one", "two")
		must.True(t, branches.Contains("one"))
		must.True(t, branches.Contains("two"))
		must.False(t, branches.Contains("three"))
	})

	t.Run("Hoist", func(t *testing.T) {
		t.Parallel()
		t.Run("haystack contains needle", func(t *testing.T) {
			t.Parallel()
			branches := gitdomain.NewLocalBranchNames("one", "two", "three")
			have := branches.Hoist(gitdomain.NewLocalBranchName("two"))
			want := gitdomain.NewLocalBranchNames("two", "one", "three")
			must.Eq(t, want, have)
		})
		t.Run("haystack does not contain needle", func(t *testing.T) {
			t.Parallel()
			branches := gitdomain.NewLocalBranchNames("one", "two", "three")
			have := branches.Hoist(gitdomain.NewLocalBranchName("zonk"))
			want := gitdomain.NewLocalBranchNames("one", "two", "three")
			must.Eq(t, want, have)
		})
	})

	t.Run("NewLocalBranchNames and Strings", func(t *testing.T) {
		t.Parallel()
		t.Run("with value", func(t *testing.T) {
			t.Parallel()
			branches := gitdomain.NewLocalBranchNames("one", "two", "three")
			want := []string{"one", "two", "three"}
			must.Eq(t, want, branches.Strings())
		})
		t.Run("no branch names", func(t *testing.T) {
			t.Parallel()
			branches := gitdomain.NewLocalBranchNames()
			must.Eq(t, 0, len(branches))
		})
	})

	t.Run("NewLocalBranchNamesRef", func(t *testing.T) {
		t.Parallel()
		t.Run("no branches", func(t *testing.T) {
			t.Parallel()
			have := gitdomain.ParseLocalBranchNamesRef("")
			must.EqOp(t, 0, len(*have))
		})
		t.Run("one branch", func(t *testing.T) {
			t.Parallel()
			have := gitdomain.ParseLocalBranchNamesRef("one")
			want := []string{"one"}
			must.Eq(t, want, have.Strings())
		})
		t.Run("multiple branches", func(t *testing.T) {
			t.Parallel()
			have := gitdomain.ParseLocalBranchNamesRef("one two three")
			want := []string{"one", "two", "three"}
			must.Eq(t, want, have.Strings())
		})
	})

	t.Run("Remove", func(t *testing.T) {
		t.Parallel()
		t.Run("the element to remove exist in the list", func(t *testing.T) {
			t.Parallel()
			branches := gitdomain.NewLocalBranchNames("one", "two", "three")
			have := branches.Remove("two")
			want := gitdomain.NewLocalBranchNames("one", "three")
			must.Eq(t, want, have)
		})
		t.Run("the element to remove does not exist in the list", func(t *testing.T) {
			t.Parallel()
			branches := gitdomain.NewLocalBranchNames("one", "two")
			have := branches.Remove("zonk")
			want := gitdomain.NewLocalBranchNames("one", "two")
			must.Eq(t, want, have)
		})
		t.Run("removing multiple elements", func(t *testing.T) {
			t.Parallel()
			branches := gitdomain.NewLocalBranchNames("one", "two", "three", "four")
			have := branches.Remove("two", "four", "zonk")
			want := gitdomain.NewLocalBranchNames("one", "three")
			must.Eq(t, want, have)
		})
	})

	t.Run("RemoveMarkers", func(t *testing.T) {
		t.Parallel()
		branches := gitdomain.NewLocalBranchNames("one", "+ two")
		have := branches.RemoveWorktreeMarkers()
		want := gitdomain.NewLocalBranchNames("one", "two")
		must.Eq(t, want, have)
	})

	t.Run("TrackingBranch", func(t *testing.T) {
		t.Parallel()
		branch := gitdomain.NewLocalBranchName("branch")
		have := branch.TrackingBranch()
		want := gitdomain.NewRemoteBranchName("origin/branch")
		must.EqOp(t, want, have)
	})

	t.Run("Sort", func(t *testing.T) {
		t.Parallel()
		branches := gitdomain.NewLocalBranchNames("one", "two", "three")
		want := []string{"one", "three", "two"}
		branches.Sort()
		must.Eq(t, want, branches.Strings())
	})
}
