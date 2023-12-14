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
			Global: gitconfig.SingleCache{
				configdomain.KeyOffline: "1",
			},
			Local: gitconfig.SingleCache{
				configdomain.KeyMainBranch: "main",
			},
		}
		clone := original.Clone()
		clone.Global[configdomain.KeyOffline] = "0"
		clone.Local[configdomain.KeyMainBranch] = "dev"
		must.EqOp(t, "1", original.Global[configdomain.KeyOffline])
		must.EqOp(t, "main", original.Local[configdomain.KeyMainBranch])
	})
}
