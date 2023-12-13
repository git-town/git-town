package configdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/shoenig/test/must"
)

func TestOffline(t *testing.T) {
	t.Parallel()

	t.Run("Online", func(t *testing.T) {
		t.Parallel()
		t.Run("is offline", func(t *testing.T) {
			t.Parallel()
			offline := configdomain.Offline(true)
			have := offline.ToOnline()
			want := configdomain.Online(false)
			must.EqOp(t, want, have)
		})
	})
}
