package data_test

import (
	"testing"

	"github.com/git-town/git-town/tools/release_stats/data"
	"github.com/shoenig/test/must"
)

func TestUsers(t *testing.T) {
	t.Parallel()
	t.Run("AddUser", func(t *testing.T) {
		t.Parallel()
		users := data.NewUsers()
		users.AddUser("one")
		users.AddUser("one")
		users.AddUser("two")
		have := users.Users()
		want := []string{"one", "two"}
		must.Eq(t, want, have)
	})
	t.Run("AddUsers", func(t *testing.T) {
		t.Parallel()
		allUsers := data.NewUsers()
		allUsers.AddUser("alpha")
		otherUsers := data.NewUsers()
		otherUsers.AddUser("beta1")
		otherUsers.AddUser("beta2")
		allUsers.AddUsers(otherUsers)
		have := allUsers.Users()
		want := []string{"alpha", "beta1", "beta2"}
		must.Eq(t, want, have)
	})
}
