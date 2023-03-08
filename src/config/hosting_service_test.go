package config_test

import (
	"testing"

	"github.com/git-town/git-town/v7/src/config"
	"github.com/stretchr/testify/assert"
)

func TestNewHostingService(t *testing.T) {
	t.Parallel()
	t.Run("valid content", func(t *testing.T) {
		t.Parallel()
		tests := map[string]config.HostingService{
			"bitbucket": config.HostingServiceBitbucket,
			"github":    config.HostingServiceGitHub,
			"gitlab":    config.HostingServiceGitLab,
			"gitea":     config.HostingServiceGitea,
			"":          config.HostingServiceNone,
		}
		for give, want := range tests {
			have, err := config.NewHostingService(give)
			assert.Nil(t, err)
			assert.Equal(t, want, have)
		}
	})

	t.Run("case insensitive", func(t *testing.T) {
		t.Parallel()
		for _, give := range []string{"github", "GitHub", "GITHUB"} {
			have, err := config.NewHostingService(give)
			assert.Nil(t, err)
			assert.Equal(t, config.HostingServiceGitHub, have)
		}
	})

	t.Run("invalid content", func(t *testing.T) {
		t.Parallel()
		_, err := config.NewHostingService("zonk")
		assert.Error(t, err)
	})
}
