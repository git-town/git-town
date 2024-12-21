package list_test

import (
	"testing"

	"github.com/git-town/git-town/v17/internal/cli/dialog/components/list"
	"github.com/git-town/git-town/v17/internal/config/configdomain"
	. "github.com/git-town/git-town/v17/pkg/prelude"
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
			entries := list.Entries[configdomain.HostingPlatform]{
				{
					Data: configdomain.HostingPlatformGitHub,
					Text: "Github",
				},
				{
					Data: configdomain.HostingPlatformGitLab,
					Text: "GitLab",
				},
			}
			have := entries.IndexOf(configdomain.HostingPlatformGitLab)
			want := 1
			must.EqOp(t, want, have)
		})
		t.Run("does not work with options", func(t *testing.T) {
			t.Parallel()
			entries := list.Entries[Option[configdomain.HostingPlatform]]{
				{
					Data: None[configdomain.HostingPlatform](),
					Text: "auto-detect",
				},
				{
					Data: Some(configdomain.HostingPlatformGitHub),
					Text: "Github",
				},
				{
					Data: Some(configdomain.HostingPlatformGitLab),
					Text: "GitLab",
				},
			}
			have := entries.IndexOf(Some(configdomain.HostingPlatformGitHub))
			want := 0 // this should be 1 if comparing options would work
			must.EqOp(t, want, have)
			have = entries.IndexOfFunc(Some(configdomain.HostingPlatformGitHub), func(a, b Option[configdomain.HostingPlatform]) bool {
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
			entries := list.Entries[Option[configdomain.HostingPlatform]]{
				{
					Data: None[configdomain.HostingPlatform](),
					Text: "auto-detect",
				},
				{
					Data: Some(configdomain.HostingPlatformGitHub),
					Text: "Github",
				},
				{
					Data: Some(configdomain.HostingPlatformGitLab),
					Text: "GitLab",
				},
			}
			have := entries.IndexOfFunc(Some(configdomain.HostingPlatformGitHub), func(a, b Option[configdomain.HostingPlatform]) bool {
				return a.Equal(b)
			})
			want := 1
			must.EqOp(t, want, have)
		})
	})
}
