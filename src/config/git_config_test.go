package config_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/stretchr/testify/assert"
)

func TestGitConfigCache(t *testing.T) {
	t.Parallel()
	t.Run("Clone", func(t *testing.T) {
		t.Parallel()
		alpha := config.Key{"alpha"}
		beta := config.Key{"beta"}
		original := config.GitConfigCache{
			alpha: "A",
			beta:  "B",
		}
		cloned := original.Clone()
		cloned[alpha] = "new A"
		cloned[beta] = "new B"
		assert.Equal(t, "A", original[alpha])
		assert.Equal(t, "B", original[beta])
	})
}
