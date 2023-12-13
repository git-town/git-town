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
		t.Run("true", func(t *testing.T) {
			t.Parallel()
			hook := domain.PushHook(true)
			have := hook.Negate()
			want := domain.PushHook(false)
			must.EqOp(t, want, have)
		})
		t.Run("false", func(t *testing.T) {
			t.Parallel()
			hook := domain.PushHook(false)
			have := hook.Negate()
			want := domain.PushHook(true)
			must.EqOp(t, want, have)
		})
	})
}
