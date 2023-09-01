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

	t.Run("KeysMatching", func(t *testing.T) {
		t.Parallel()
		cache := config.GitConfigCache{
			config.Key{"key1"}:  "A",
			config.Key{"key2"}:  "B",
			config.Key{"other"}: "other",
		}
		have := cache.KeysMatching("key")
		want := []config.Key{
			{"key1"},
			{"key2"},
		}
		assert.Equal(t, want, have)
	})
}
