package configdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v17/internal/config/configdomain"
	. "github.com/git-town/git-town/v17/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestNewHostingPlatform(t *testing.T) {
	t.Parallel()

	t.Run("valid content", func(t *testing.T) {
		t.Parallel()
		tests := map[string]Option[configdomain.HostingPlatform]{
			"":                     None[configdomain.HostingPlatform](),
			"bitbucket":            Some(configdomain.HostingPlatformBitbucket),
			"BitBucket":            Some(configdomain.HostingPlatformBitbucket),
			"BITBUCKET":            Some(configdomain.HostingPlatformBitbucket),
			"bitbucket-datacenter": Some(configdomain.HostingPlatformBitbucketDatacenter),
			"BitBucket-Datacenter": Some(configdomain.HostingPlatformBitbucketDatacenter),
			"BITBUCKET-DATACENTER": Some(configdomain.HostingPlatformBitbucketDatacenter),
			"github":               Some(configdomain.HostingPlatformGitHub),
			"GitHub":               Some(configdomain.HostingPlatformGitHub),
			"gitlab":               Some(configdomain.HostingPlatformGitLab),
			"GitLab":               Some(configdomain.HostingPlatformGitLab),
			"gitea":                Some(configdomain.HostingPlatformGitea),
			"Gitea":                Some(configdomain.HostingPlatformGitea),
		}
		for give, want := range tests {
			have, err := configdomain.ParseHostingPlatform(give)
			must.NoError(t, err)
			must.Eq(t, want, have)
		}
	})

	t.Run("invalid content", func(t *testing.T) {
		t.Parallel()
		_, err := configdomain.ParseHostingPlatform("zonk")
		must.Error(t, err)
	})
}
