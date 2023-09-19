package persistence_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/persistence"
	"github.com/stretchr/testify/assert"
)

func TestSanitizePath(t *testing.T) {
	t.Parallel()
	t.Run("SanitizePath", func(t *testing.T) {
		t.Parallel()
		tests := map[string]string{
			"/home/user/development/git-town":        "home-user-development-git-town",
			"c:\\Users\\user\\development\\git-town": "c-users-user-development-git-town",
		}
		for give, want := range tests {
			have := persistence.SanitizePath(give)
			assert.Equal(t, want, have)
		}
	})
}
