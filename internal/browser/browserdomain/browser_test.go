package browserdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/browser/browserdomain"
	"github.com/shoenig/test/must"
)

func TestUseBrowser(t *testing.T) {
	t.Parallel()

	t.Run("ParseBrowserHas", func(t *testing.T) {
		t.Parallel()

		t.Run("not configured", func(t *testing.T) {
			t.Parallel()
			haveExecutable, haveEnabled, err := browserdomain.ParseBrowserHas("", false)
			must.NoError(t, err)
			must.True(t, haveExecutable.IsNone())
			must.True(t, haveEnabled.IsNone())
		})

		t.Run("set to empty string", func(t *testing.T) {
			t.Parallel()
			haveExecutable, haveEnabled, err := browserdomain.ParseBrowserHas("", false)
			must.NoError(t, err)
			must.True(t, haveExecutable.IsNone())
			must.True(t, haveEnabled.EqualSome(browserdomain.BrowserEnabled(false)))
		})

		t.Run("set to '(none)'", func(t *testing.T) {
			t.Parallel()
			haveExecutable, haveEnabled, err := browserdomain.ParseBrowserHas("(none)", false)
			must.NoError(t, err)
			must.True(t, haveExecutable.IsNone())
			must.True(t, haveEnabled.EqualSome(browserdomain.BrowserEnabled(false)))
		})

		t.Run("set to an actual browser executable", func(t *testing.T) {
			t.Parallel()
			haveExecutable, haveEnabled, err := browserdomain.ParseBrowserHas("firefox", false)
			must.NoError(t, err)
			must.True(t, haveExecutable.EqualSome(browserdomain.BrowserExecutable("firefox")))
			must.True(t, haveEnabled.EqualSome(browserdomain.BrowserEnabled(true)))
		})
	})
}
