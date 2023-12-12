package domain_test

import (
	"testing"

	"github.com/git-town/git-town/v11/src/domain"
	"github.com/shoenig/test/must"
)

func TestLocalBranchNames(t *testing.T) {
	t.Parallel()

	t.Run("AtRemote", func(t *testing.T) {
		t.Parallel()
		branch := domain.NewLocalBranchName("branch")
		have := branch.AtRemote(domain.OriginRemote)
		want := domain.NewRemoteBranchName("origin/branch")
		must.EqOp(t, want, have)
	})

	t.Run("Hoist", func(t *testing.T) {
		t.Parallel()
		t.Run("haystack contains needle", func(t *testing.T) {
			t.Parallel()
			branches := domain.NewLocalBranchNames("one", "two", "three")
			have := branches.Hoist(domain.NewLocalBranchName("two"))
			want := domain.NewLocalBranchNames("two", "one", "three")
			must.Eq(t, want, have)
		})
		t.Run("haystack does not contain needle", func(t *testing.T) {
			t.Parallel()
			branches := domain.NewLocalBranchNames("one", "two", "three")
			have := branches.Hoist(domain.NewLocalBranchName("zonk"))
			want := domain.NewLocalBranchNames("one", "two", "three")
			must.Eq(t, want, have)
		})
	})

	t.Run("NewLocalBranchNames and Strings", func(t *testing.T) {
		t.Parallel()
		branches := domain.NewLocalBranchNames("one", "two", "three")
		want := []string{"one", "two", "three"}
		must.Eq(t, want, branches.Strings())
	})

	t.Run("Remove", func(t *testing.T) {
		t.Parallel()
		t.Run("the element to remove exist in the list", func(t *testing.T) {
			t.Parallel()
			branches := domain.NewLocalBranchNames("one", "two", "three")
			have := branches.Remove("two")
			want := domain.NewLocalBranchNames("one", "three")
			must.Eq(t, want, have)
		})
		t.Run("the element to remove does not exist in the list", func(t *testing.T) {
			t.Parallel()
			branches := domain.NewLocalBranchNames("one", "two")
			have := branches.Remove("zonk")
			want := domain.NewLocalBranchNames("one", "two")
			must.Eq(t, want, have)
		})
	})

	t.Run("RemoveMarkers", func(t *testing.T) {
		t.Parallel()
		branches := domain.NewLocalBranchNames("one", "+ two")
		have := branches.RemoveWorkspaceMarkers()
		want := domain.NewLocalBranchNames("one", "two")
		must.Eq(t, want, have)
	})

	t.Run("TrackingBranch", func(t *testing.T) {
		t.Parallel()
		branch := domain.NewLocalBranchName("branch")
		have := branch.TrackingBranch()
		want := domain.NewRemoteBranchName("origin/branch")
		must.EqOp(t, want, have)
	})

	t.Run("Sort", func(t *testing.T) {
		t.Parallel()
		branches := domain.NewLocalBranchNames("one", "two", "three")
		want := []string{"one", "three", "two"}
		branches.Sort()
		must.Eq(t, want, branches.Strings())
	})
}
