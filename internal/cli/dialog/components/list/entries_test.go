package list_test

import (
	"testing"

	"github.com/git-town/git-town/v21/internal/cli/dialog/components/list"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestEntries(t *testing.T) {
	t.Parallel()

	t.Run("AllDisabled", func(t *testing.T) {
		t.Parallel()
		t.Run("all entries are disabled", func(t *testing.T) {
			t.Parallel()
			entries := list.Entries[configdomain.HostingOriginHostname]{
				{Disabled: true},
				{Disabled: true},
			}
			must.True(t, entries.AllDisabled())
		})
		t.Run("some entries are enabled", func(t *testing.T) {
			t.Parallel()
			entries := list.Entries[configdomain.HostingOriginHostname]{
				{Disabled: true},
				{Disabled: true},
				{Disabled: false},
			}
			must.False(t, entries.AllDisabled())
		})
	})

	t.Run("IndexOf", func(t *testing.T) {
		t.Parallel()
		t.Run("works with correctly comparable types", func(t *testing.T) {
			t.Parallel()
			entries := list.Entries[forgedomain.ForgeType]{
				{
					Data: forgedomain.ForgeTypeGitHub,
					Text: "Github",
				},
				{
					Data: forgedomain.ForgeTypeGitLab,
					Text: "GitLab",
				},
			}
			have := entries.IndexOf(forgedomain.ForgeTypeGitLab)
			want := 1
			must.EqOp(t, want, have)
		})
		t.Run("does not work with options", func(t *testing.T) {
			t.Parallel()
			entries := list.Entries[Option[forgedomain.ForgeType]]{
				{
					Data: None[forgedomain.ForgeType](),
					Text: "auto-detect",
				},
				{
					Data: Some(forgedomain.ForgeTypeGitHub),
					Text: "Github",
				},
				{
					Data: Some(forgedomain.ForgeTypeGitLab),
					Text: "GitLab",
				},
			}
			have := entries.IndexOf(Some(forgedomain.ForgeTypeGitHub))
			want := 0 // this should be 1 if comparing options would work
			must.EqOp(t, want, have)
			have = entries.IndexOfFunc(Some(forgedomain.ForgeTypeGitHub), func(a, b Option[forgedomain.ForgeType]) bool {
				return a.Equal(b)
			})
			want = 1
			must.EqOp(t, want, have)
		})
	})

	t.Run("IndexOfFunc", func(t *testing.T) {
		t.Parallel()
		t.Run("works with options", func(t *testing.T) {
			t.Parallel()
			entries := list.Entries[Option[forgedomain.ForgeType]]{
				{
					Data: None[forgedomain.ForgeType](),
					Text: "auto-detect",
				},
				{
					Data: Some(forgedomain.ForgeTypeGitHub),
					Text: "Github",
				},
				{
					Data: Some(forgedomain.ForgeTypeGitLab),
					Text: "GitLab",
				},
			}
			have := entries.IndexOfFunc(Some(forgedomain.ForgeTypeGitHub), func(a, b Option[forgedomain.ForgeType]) bool {
				return a.Equal(b)
			})
			want := 1
			must.EqOp(t, want, have)
		})
	})
}
