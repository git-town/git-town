package list_test

import (
	"testing"

	"github.com/git-town/git-town/v18/internal/cli/dialog/components/list"
	"github.com/git-town/git-town/v18/internal/config/configdomain"
	. "github.com/git-town/git-town/v18/pkg/prelude"
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
			entries := list.Entries[configdomain.ForgeType]{
				{
					Data: configdomain.ForgeTypeGitHub,
					Text: "Github",
				},
				{
					Data: configdomain.ForgeTypeGitLab,
					Text: "GitLab",
				},
			}
			have := entries.IndexOf(configdomain.ForgeTypeGitLab)
			want := 1
			must.EqOp(t, want, have)
		})
		t.Run("does not work with options", func(t *testing.T) {
			t.Parallel()
			entries := list.Entries[Option[configdomain.ForgeType]]{
				{
					Data: None[configdomain.ForgeType](),
					Text: "auto-detect",
				},
				{
					Data: Some(configdomain.ForgeTypeGitHub),
					Text: "Github",
				},
				{
					Data: Some(configdomain.ForgeTypeGitLab),
					Text: "GitLab",
				},
			}
			have := entries.IndexOf(Some(configdomain.ForgeTypeGitHub))
			want := 0 // this should be 1 if comparing options would work
			must.EqOp(t, want, have)
			have = entries.IndexOfFunc(Some(configdomain.ForgeTypeGitHub), func(a, b Option[configdomain.ForgeType]) bool {
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
			entries := list.Entries[Option[configdomain.ForgeType]]{
				{
					Data: None[configdomain.ForgeType](),
					Text: "auto-detect",
				},
				{
					Data: Some(configdomain.ForgeTypeGitHub),
					Text: "Github",
				},
				{
					Data: Some(configdomain.ForgeTypeGitLab),
					Text: "GitLab",
				},
			}
			have := entries.IndexOfFunc(Some(configdomain.ForgeTypeGitHub), func(a, b Option[configdomain.ForgeType]) bool {
				return a.Equal(b)
			})
			want := 1
			must.EqOp(t, want, have)
		})
	})
}
