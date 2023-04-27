package asserts_test

import (
	"errors"
	"testing"

	"github.com/git-town/git-town/v8/test/asserts"
)

func TestOk(t *testing.T) {
	t.Parallel()
	t.Run("giving nil does not panic", func(t *testing.T) {
		t.Parallel()
		asserts.Ok(nil)
	})
	t.Run("giving error panics", func(t *testing.T) {
		t.Parallel()
		asserts.Panics(t, func() {
			asserts.Ok(errors.New("BOOM"))
		})
	})
}
