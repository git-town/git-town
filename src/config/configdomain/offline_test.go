package configdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v13/src/config/configdomain"
	"github.com/shoenig/test/must"
)

func TestOffline(t *testing.T) {
	t.Parallel()

	t.Run("ToOnline", func(t *testing.T) {
		t.Parallel()
		tests := map[bool]bool{
			true:  false,
			false: true,
		}
		for give, wantBool := range tests {
			offline := configdomain.Offline(give)
			have := offline.ToOnline()
			want := configdomain.Online(wantBool)
			must.EqOp(t, want, have)
		}
	})
}
