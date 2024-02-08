package main_test

import (
	"testing"

	listContributors "github.com/git-town/git-town/tools/list_contributors"
	"github.com/shoenig/test/must"
)

func TestUserCollector(t *testing.T) {
	t.Parallel()
	uc := listContributors.UserCollector{}
	uc.AddUser("one")
	uc.AddUser("one")
	uc.AddUser("two")
	have := uc.Users()
	want := []string{"one", "two"}
	must.Eq(t, want, have)
}
