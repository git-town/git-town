package gitconfig_test

import (
	"testing"

	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/config/gitconfig"
	"github.com/shoenig/test/must"
)

func TestGitConfig(t *testing.T) {
	t.Parallel()

	t.Run("Clone", func(t *testing.T) {
		t.Parallel()
		original := gitconfig.FullCache{
			GlobalCache: gitconfig.SingleCache{
				configdomain.KeyOffline: "1",
			},
			LocalCache: gitconfig.SingleCache{
				configdomain.KeyMainBranch: "main",
			},
		}
		clone := original.Clone()
		clone.GlobalCache[configdomain.KeyOffline] = "0"
		clone.LocalCache[configdomain.KeyMainBranch] = "dev"
		must.EqOp(t, "1", original.GlobalCache[configdomain.KeyOffline])
		must.EqOp(t, "main", original.LocalCache[configdomain.KeyMainBranch])
	})
}
