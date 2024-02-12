package data_test

import (
	"testing"

	"github.com/git-town/git-town/tools/stats_release/data"
	"github.com/shoenig/test/must"
)

func TestUsers(t *testing.T) {
	t.Parallel()

	t.Run("provides contributors sorted by contribution count", func(t *testing.T) {
		t.Parallel()
		counter := data.NewContributionCounter()
		counter.AddUser("one")
		counter.AddUser("two")
		counter.AddUser("two")
		have := counter.Contributors()
		want := []data.Contributor{
			{
				Username:          "two",
				ContributionCount: 2,
			},
			{
				Username:          "one",
				ContributionCount: 1,
			},
		}
		must.Eq(t, want, have)
	})

	t.Run("sorts contributors with the same contribution amount alphabetically", func(t *testing.T) {
		t.Parallel()
		counter := data.NewContributionCounter()
		counter.AddUser("gamma")
		counter.AddUser("beta")
		counter.AddUser("alpha")
		have := counter.Contributors()
		want := []data.Contributor{
			{
				Username:          "alpha",
				ContributionCount: 1,
			},
			{
				Username:          "beta",
				ContributionCount: 1,
			},
			{
				Username:          "gamma",
				ContributionCount: 1,
			},
		}
		must.Eq(t, want, have)
	})

	t.Run("allows adding multiple contributions efficiently", func(t *testing.T) {
		t.Parallel()
		allUsers := data.NewContributionCounter()
		allUsers.AddUser("alpha")
		otherUsers := data.NewContributionCounter()
		otherUsers.AddUser("alpha")
		otherUsers.AddUser("beta2")
		otherUsers.AddUser("beta1")
		allUsers.AddUsers(otherUsers)
		have := allUsers.Contributors()
		want := []data.Contributor{
			{
				Username:          "alpha",
				ContributionCount: 2,
			},
			{
				Username:          "beta1",
				ContributionCount: 2,
			},
			{
				Username:          "beta2",
				ContributionCount: 2,
			},
		}
		must.Eq(t, want, have)
	})
}
