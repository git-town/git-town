package configdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v13/src/config/configdomain"
	"github.com/shoenig/test/must"
)

func TestNewHostingPlatform(t *testing.T) {
	t.Parallel()

	t.Run("valid content", func(t *testing.T) {
		t.Parallel()
		tests := map[string]configdomain.HostingPlatform{
			"bitbucket": configdomain.HostingPlatformBitbucket,
			"BitBucket": configdomain.HostingPlatformBitbucket,
			"github":    configdomain.HostingPlatformGitHub,
			"GitHub":    configdomain.HostingPlatformGitHub,
			"gitlab":    configdomain.HostingPlatformGitLab,
			"GitLab":    configdomain.HostingPlatformGitLab,
			"gitea":     configdomain.HostingPlatformGitea,
			"Gitea":     configdomain.HostingPlatformGitea,
			"":          configdomain.HostingPlatformNone,
		}
		for give, want := range tests {
			have, err := configdomain.NewHostingPlatform(give)
			must.NoError(t, err)
			must.EqOp(t, want, have)
		}
	})

	t.Run("case insensitive", func(t *testing.T) {
		t.Parallel()
		for _, give := range []string{"github", "GitHub", "GITHUB"} {
			have, err := configdomain.NewHostingPlatform(give)
			must.NoError(t, err)
			must.EqOp(t, configdomain.HostingPlatformGitHub, have)
		}
	})

	t.Run("invalid content", func(t *testing.T) {
		t.Parallel()
		_, err := configdomain.NewHostingPlatform("zonk")
		must.Error(t, err)
	})
}
