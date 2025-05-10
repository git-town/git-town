package configdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v20/internal/config/configdomain"
	"github.com/shoenig/test/must"
)

func TestOffline(t *testing.T) {
	t.Parallel()

	t.Run("IsOffline", func(t *testing.T) {
		t.Parallel()
		tests := map[bool]bool{
			true:  true,
			false: false,
		}
		for give, want := range tests {
			offline := configdomain.Offline(give)
			have := offline.IsOffline()
			must.EqOp(t, want, have)
		}
	})

	t.Run("IsOffline", func(t *testing.T) {
		t.Parallel()
		tests := map[bool]bool{
			true:  false,
			false: true,
		}
		for give, want := range tests {
			offline := configdomain.Offline(give)
			have := offline.IsOnline()
			must.EqOp(t, want, have)
		}
	})
}
