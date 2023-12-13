package domain_test

import (
	"testing"

	"github.com/git-town/git-town/v11/src/domain"
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
			hook := domain.PushHook(giveBool)
			have := hook.Negate()
			want := domain.PushHook(wantBool)
			must.EqOp(t, want, have)
		}
	})
}
