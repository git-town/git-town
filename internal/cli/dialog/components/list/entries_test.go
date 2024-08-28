package list_test

import (
	"testing"

	"github.com/git-town/git-town/v15/internal/cli/dialog/components/list"
	"github.com/git-town/git-town/v15/internal/config/configdomain"
	"github.com/shoenig/test/must"
)

func TestEntries(t *testing.T) {
	t.Parallel()

	t.Run("AllDisabled", func(t *testing.T) {
		t.Parallel()
		t.Run("all entries are disabled", func(t *testing.T) {
			t.Parallel()
			entries := list.Entries[configdomain.HostingOriginHostname]{
				{Enabled: false},
				{Enabled: false},
			}
			must.True(t, entries.AllDisabled())
		})
		t.Run("some entries are enabled", func(t *testing.T) {
			t.Parallel()
			entries := list.Entries[configdomain.HostingOriginHostname]{
				{Enabled: false},
				{Enabled: false},
				{Enabled: true},
			}
			must.False(t, entries.AllDisabled())
		})
	})
}
