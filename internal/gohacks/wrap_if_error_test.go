package gohacks_test

import (
	"errors"
	"testing"

	"github.com/git-town/git-town/v22/internal/gohacks"
	"github.com/shoenig/test/must"
)

func TestWrapIfError(t *testing.T) {
	t.Parallel()

	t.Run("error", func(t *testing.T) {
		t.Parallel()
		err := errors.New("my error")
		have := gohacks.WrapIfError(err, "encountered an error in file %q: %v", "my_file.txt").Error()
		want := `encountered an error in file "my_file.txt": my error`
		must.EqOp(t, want, have)
	})

	t.Run("no error", func(t *testing.T) {
		t.Parallel()
		have := gohacks.WrapIfError(nil, "encountered error: %v")
		must.Nil(t, have)
	})

	t.Run("only error", func(t *testing.T) {
		t.Parallel()
		err := errors.New("my error")
		have := gohacks.WrapIfError(err, "encountered error: %v").Error()
		want := `encountered error: my error`
		must.EqOp(t, want, have)
	})
}
