package main_test

import (
	"testing"

	listContributors "github.com/git-town/git-town/tools/list_contributors"
	"github.com/shoenig/test/must"
)

func TestUserCollector(t *testing.T) {
	t.Parallel()
	t.Run("AddUser", func(t *testing.T) {
		users := listContributors.NewUsers()
		users.AddUser("one")
		users.AddUser("one")
		users.AddUser("two")
		have := users.Users()
		want := []string{"one", "two"}
		must.Eq(t, want, have)
	})
	t.Run("AddUsers", func(t *testing.T) {
		totalUsers := listContributors.NewUsers()
		totalUsers.AddUser("alpha")
		issueUsers := listContributors.NewUsers()
		issueUsers.AddUser("beta1")
		issueUsers.AddUser("beta2")
		totalUsers.AddUsers(issueUsers)
		have := totalUsers.Users()
		want := []string{"alpha", "beta1", "beta2"}
		must.Eq(t, want, have)
	})
}
