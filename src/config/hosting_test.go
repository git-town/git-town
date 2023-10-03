package config_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/shoenig/test"
)

func TestNewHostingService(t *testing.T) {
	t.Parallel()

	t.Run("valid content", func(t *testing.T) {
		t.Parallel()
		tests := map[string]config.Hosting{
			"bitbucket": config.HostingBitbucket,
			"github":    config.HostingGitHub,
			"gitlab":    config.HostingGitLab,
			"gitea":     config.HostingGitea,
			"":          config.HostingNone,
		}
		for give, want := range tests {
			have, err := config.NewHosting(give)
			test.NoError(t, err)
			test.EqOp(t, want, have)
		}
	})

	t.Run("case insensitive", func(t *testing.T) {
		t.Parallel()
		for _, give := range []string{"github", "GitHub", "GITHUB"} {
			have, err := config.NewHosting(give)
			test.NoError(t, err)
			test.EqOp(t, config.HostingGitHub, have)
		}
	})

	t.Run("invalid content", func(t *testing.T) {
		t.Parallel()
		_, err := config.NewHosting("zonk")
		test.Error(t, err)
	})
}
