package domain_test

import (
	"testing"

	"github.com/git-town/git-town/v10/src/domain"
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

	t.Run("NewLocalBranchNames and Strings", func(t *testing.T) {
		t.Parallel()
		branches := domain.NewLocalBranchNames("one", "two", "three")
		want := []string{"one", "two", "three"}
		must.Eq(t, want, branches.Strings())
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
