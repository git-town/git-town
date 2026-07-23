package browserdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v24/internal/browser/browserdomain"
	. "github.com/git-town/git-town/v24/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestUseBrowser(t *testing.T) {
	t.Parallel()

	t.Run("ParseBrowserOpt", func(t *testing.T) {
		t.Parallel()

		t.Run("not configured", func(t *testing.T) {
			t.Parallel()
			haveExecutable, haveEnabled, err := browserdomain.ParseBrowserOpt(None[string]())
			must.NoError(t, err)
			must.True(t, haveExecutable.IsNone())
			must.True(t, haveEnabled.IsNone())
		})

		t.Run("set to empty string", func(t *testing.T) {
			t.Parallel()
			haveExecutable, haveEnabled, err := browserdomain.ParseBrowserOpt(Some(""))
			must.NoError(t, err)
			must.True(t, haveExecutable.IsNone())
			must.True(t, haveEnabled.EqualSome(browserdomain.BrowserEnabled(false)))
		})

		t.Run("set to '(none)'", func(t *testing.T) {
			t.Parallel()
			haveExecutable, haveEnabled, err := browserdomain.ParseBrowserOpt(Some("(none)"))
			must.NoError(t, err)
			must.True(t, haveExecutable.IsNone())
			must.True(t, haveEnabled.EqualSome(browserdomain.BrowserEnabled(false)))
		})

		t.Run("set to an actual browser executable", func(t *testing.T) {
			t.Parallel()
			haveExecutable, haveEnabled, err := browserdomain.ParseBrowserOpt(Some("firefox"))
			must.NoError(t, err)
			must.True(t, haveExecutable.EqualSome(browserdomain.BrowserExecutable("firefox")))
			must.True(t, haveEnabled.IsNone())
		})
	})
}
