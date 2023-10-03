package config_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/shoenig/test"
)

func TestGitConfig(t *testing.T) {
	t.Parallel()

	t.Run("Clone", func(t *testing.T) {
		t.Parallel()
		original := config.GitConfig{
			Global: config.GitConfigCache{
				config.KeyOffline: "1",
			},
			Local: config.GitConfigCache{
				config.KeyMainBranch: "main",
			},
		}
		clone := original.Clone()
		clone.Global[config.KeyOffline] = "0"
		clone.Local[config.KeyMainBranch] = "dev"
		test.EqOp(t, "1", original.Global[config.KeyOffline])
		test.EqOp(t, "main", original.Local[config.KeyMainBranch])
	})
}
