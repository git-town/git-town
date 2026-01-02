package list_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents/list"
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
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

	t.Run("FirstEnabled", func(t *testing.T) {
		t.Parallel()
		t.Run("first entry is disabled", func(t *testing.T) {
			t.Parallel()
			entries := list.Entries[configdomain.HostingOriginHostname]{
				{Disabled: true},
				{Disabled: false},
			}
			have := entries.FirstEnabled()
			must.EqOp(t, 1, have)
		})
		t.Run("second entry is disabled", func(t *testing.T) {
			t.Parallel()
			entries := list.Entries[configdomain.HostingOriginHostname]{
				{Disabled: false},
				{Disabled: true},
			}
			have := entries.FirstEnabled()
			must.EqOp(t, 0, have)
		})
		t.Run("all entries are disabled", func(t *testing.T) {
			t.Parallel()
			entries := list.Entries[configdomain.HostingOriginHostname]{
				{Disabled: true},
				{Disabled: true},
			}
			have := entries.FirstEnabled()
			must.EqOp(t, 0, have)
		})
		t.Run("no entries", func(t *testing.T) {
			t.Parallel()
			entries := list.Entries[configdomain.HostingOriginHostname]{}
			have := entries.FirstEnabled()
			must.EqOp(t, 0, have)
		})
	})

	t.Run("IndexOf", func(t *testing.T) {
		t.Parallel()
		t.Run("works with comparable types", func(t *testing.T) {
			t.Parallel()
			entries := list.Entries[forgedomain.ForgeType]{
				{
					Data: forgedomain.ForgeTypeGithub,
					Text: "GitHub",
				},
				{
					Data: forgedomain.ForgeTypeGitlab,
					Text: "GitLab",
				},
			}
			have := entries.IndexOf(forgedomain.ForgeTypeGitlab)
			want := 1
			must.EqOp(t, want, have)
		})
		t.Run("works with options", func(t *testing.T) {
			t.Parallel()
			entries := list.Entries[Option[forgedomain.ForgeType]]{
				{
					Data: None[forgedomain.ForgeType](),
					Text: "auto-detect",
				},
				{
					Data: Some(forgedomain.ForgeTypeGithub),
					Text: "GitHub",
				},
				{
					Data: Some(forgedomain.ForgeTypeGitlab),
					Text: "GitLab",
				},
			}
			have := entries.IndexOf(Some(forgedomain.ForgeTypeGithub))
			want := 1
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
					Data: Some(forgedomain.ForgeTypeGithub),
					Text: "GitHub",
				},
				{
					Data: Some(forgedomain.ForgeTypeGitlab),
					Text: "GitLab",
				},
			}
			have := entries.IndexOfFunc(Some(forgedomain.ForgeTypeGithub), func(a, b Option[forgedomain.ForgeType]) bool {
				return a.Equal(b)
			})
			want := 1
			must.EqOp(t, want, have)
		})
	})
}
