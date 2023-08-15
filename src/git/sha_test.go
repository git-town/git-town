package git_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/git"
	"github.com/stretchr/testify/assert"
)

func ensurePanic(t *testing.T) {
	if r := recover(); r == nil {
		t.Errorf("The code did not panic")
	}
}

func TestSHA(t *testing.T) {
	t.Run("NewSHA", func(t *testing.T) {
		t.Run("allows hex characters", func(t *testing.T) {
			text := "1234567890abcdef"
			git.NewSHA(text) // should not panic
		})
		t.Run("does not allow empty values", func(t *testing.T) {
			defer ensurePanic(t)
			git.NewSHA("")
		})
		t.Run("does not allow non-SHA characters", func(t *testing.T) {
			defer ensurePanic(t)
			git.NewSHA("abc def")
		})
	})

	t.Run("Stringer interface", func(t *testing.T) {
		sha := git.NewSHA("abcdef")
		assert.Equal(t, "abcdef", sha.String())
	})
}
