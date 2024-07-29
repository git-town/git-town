package data_test

import (
	"testing"

	"github.com/git-town/git-town/tools/stats_release/data"
	"github.com/shoenig/test/must"
)

func TestUsers(t *testing.T) {
	t.Parallel()

	t.Run("Add", func(t *testing.T) {
		t.Parallel()
		users := data.NewUsers()
		users.Add("one")
		users.Add("one")
		users.Add("two")
		have := users.Values()
		want := []string{"one", "two"}
		must.Eq(t, want, have)
	})

	t.Run("AddUsers", func(t *testing.T) {
		t.Parallel()
		allUsers := data.NewUsers("alpha")
		otherUsers := data.NewUsers("beta1", "beta2")
		allUsers.AddUsers(otherUsers)
		have := allUsers.Values()
		want := []string{"alpha", "beta1", "beta2"}
		must.Eq(t, want, have)
	})
}
