package git_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/git"
	"github.com/stretchr/testify/assert"
)

func ensureDidPanic(t *testing.T) {
	t.Helper()
	if r := recover(); r == nil {
		t.Errorf("The code did not panic")
	}
}

func TestSHA(t *testing.T) {
	t.Parallel()
	t.Run("NewSHA", func(t *testing.T) {
		t.Parallel()
		t.Run("allows lowercase hex characters", func(t *testing.T) {
			t.Parallel()
			text := "1234567890abcdef"
			git.NewSHA(text) // should not panic
		})
		t.Run("does not allow empty values", func(t *testing.T) {
			t.Parallel()
			defer ensureDidPanic(t)
			git.NewSHA("")
		})
		t.Run("does not allow spaces", func(t *testing.T) {
			t.Parallel()
			defer ensureDidPanic(t)
			git.NewSHA("abc def")
		})
		t.Run("does not allow uppercase characters", func(t *testing.T) {
			t.Parallel()
			defer ensureDidPanic(t)
			git.NewSHA("ABCDEF")
		})
		t.Run("does not allow non-hex characters", func(t *testing.T) {
			t.Parallel()
			defer ensureDidPanic(t)
			git.NewSHA("abcdefg")
		})
	})

	t.Run("implements the Stringer interface", func(t *testing.T) {
		t.Parallel()
		sha := git.NewSHA("abcdef")
		assert.Equal(t, "abcdef", sha.String())
	})
}
