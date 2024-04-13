package configdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/shoenig/test/must"
)

func TestPushHook(t *testing.T) {
	t.Parallel()

	t.Run("Negate", func(t *testing.T) {
		t.Parallel()
		tests := map[bool]bool{
			true:  false,
			false: true,
		}
		for giveBool, wantBool := range tests {
			hook := configdomain.PushHook(giveBool)
			have := hook.Negate()
			want := configdomain.NoPushHook(wantBool)
			must.EqOp(t, want, have)
		}
	})
}
