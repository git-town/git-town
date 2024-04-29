package prelude_test

import (
	"testing"

	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/shoenig/test/must"
)

func TestSafePointer(t *testing.T) {
	t.Parallel()
	t.Run("give non-nil", func(t *testing.T) {
		t.Parallel()
		value := "hello"
		instance := NewSafePointer(&value)
		have := instance.Get()
		must.Eq(t, &value, have)
	})
}
