package configdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v18/internal/config/configdomain"
	. "github.com/git-town/git-town/v18/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestParseForgeType(t *testing.T) {
	t.Parallel()

	t.Run("acceptable content", func(t *testing.T) {
		t.Parallel()
		tests := map[string]Option[configdomain.ForgeType]{
			"":                     None[configdomain.ForgeType](),
			"bitbucket":            Some(configdomain.ForgeTypeBitbucket),
			"BitBucket":            Some(configdomain.ForgeTypeBitbucket),
			"BITBUCKET":            Some(configdomain.ForgeTypeBitbucket),
			"bitbucket-datacenter": Some(configdomain.ForgeTypeBitbucketDatacenter),
			"BitBucket-Datacenter": Some(configdomain.ForgeTypeBitbucketDatacenter),
			"BITBUCKET-DATACENTER": Some(configdomain.ForgeTypeBitbucketDatacenter),
			"github":               Some(configdomain.ForgeTypeGitHub),
			"GitHub":               Some(configdomain.ForgeTypeGitHub),
			"gitlab":               Some(configdomain.ForgeTypeGitLab),
			"GitLab":               Some(configdomain.ForgeTypeGitLab),
			"gitea":                Some(configdomain.ForgeTypeGitea),
			"Gitea":                Some(configdomain.ForgeTypeGitea),
		}
		for give, want := range tests {
			have, err := configdomain.ParseForgeType(give)
			must.NoError(t, err)
			must.Eq(t, want, have)
		}
	})

	t.Run("invalid content", func(t *testing.T) {
		t.Parallel()
		_, err := configdomain.ParseForgeType("zonk")
		must.Error(t, err)
	})
}
