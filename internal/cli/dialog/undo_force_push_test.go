package dialog_test

import (
	"testing"

	"github.com/git-town/git-town/v14/src/cli/dialog"
	"github.com/shoenig/test/must"
)

func TestBoolEntry(t *testing.T) {
	t.Parallel()

	t.Run("String", func(t *testing.T) {
		t.Parallel()
		entry := dialog.BoolEntry(true)
		have := entry.String()
		want := "true"
		must.EqOp(t, want, have)
	})
}
