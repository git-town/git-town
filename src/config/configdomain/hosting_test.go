package configdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/shoenig/test/must"
)

func TestNewHostingService(t *testing.T) {
	t.Parallel()

	t.Run("valid content", func(t *testing.T) {
		t.Parallel()
		tests := map[string]configdomain.Hosting{
			"bitbucket": configdomain.HostingBitbucket,
			"github":    configdomain.HostingGitHub,
			"gitlab":    configdomain.HostingGitLab,
			"gitea":     configdomain.HostingGitea,
			"":          configdomain.HostingNone,
		}
		for give, want := range tests {
			have, err := configdomain.NewHosting(give)
			must.NoError(t, err)
			must.EqOp(t, want, have)
		}
	})

	t.Run("case insensitive", func(t *testing.T) {
		t.Parallel()
		for _, give := range []string{"github", "GitHub", "GITHUB"} {
			have, err := configdomain.NewHosting(give)
			must.NoError(t, err)
			must.EqOp(t, configdomain.HostingGitHub, have)
		}
	})

	t.Run("invalid content", func(t *testing.T) {
		t.Parallel()
		_, err := configdomain.NewHosting("zonk")
		must.Error(t, err)
	})
}
