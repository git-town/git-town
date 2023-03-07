package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToHostingService(t *testing.T) {
	t.Parallel()
	t.Run("valid content", func(t *testing.T) {
		tests := map[string]HostingService{
			"bitbucket": HostingServiceBitbucket,
			"github":    HostingServiceGitHub,
			"gitlab":    HostingServiceGitLab,
			"gitea":     HostingServiceGitea,
			"":          HostingServiceNone,
		}
		for give, want := range tests {
			have, err := toHostingService(give)
			assert.Nil(t, err)
			assert.Equal(t, want, have)
		}
	})

	t.Run("invalid content", func(t *testing.T) {
		_, err := toHostingService("zonk")
		assert.Error(t, err)
	})
}
