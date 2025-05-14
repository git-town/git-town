package forgedomain_test

import (
	"testing"

	"github.com/git-town/git-town/v20/internal/forge/forgedomain"
	. "github.com/git-town/git-town/v20/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestParseForgeType(t *testing.T) {
	t.Parallel()

	t.Run("acceptable content", func(t *testing.T) {
		t.Parallel()
		tests := map[string]Option[forgedomain.ForgeType]{
			"":                     None[forgedomain.ForgeType](),
			"bitbucket":            Some(forgedomain.ForgeTypeBitbucket),
			"BitBucket":            Some(forgedomain.ForgeTypeBitbucket),
			"BITBUCKET":            Some(forgedomain.ForgeTypeBitbucket),
			"bitbucket-datacenter": Some(forgedomain.ForgeTypeBitbucketDatacenter),
			"BitBucket-Datacenter": Some(forgedomain.ForgeTypeBitbucketDatacenter),
			"BITBUCKET-DATACENTER": Some(forgedomain.ForgeTypeBitbucketDatacenter),
			"codeberg":             Some(forgedomain.ForgeTypeCodeberg),
			"Codeberg":             Some(forgedomain.ForgeTypeCodeberg),
			"CodeBerg":             Some(forgedomain.ForgeTypeCodeberg),
			"CODEBERG":             Some(forgedomain.ForgeTypeCodeberg),
			"github":               Some(forgedomain.ForgeTypeGitHub),
			"GitHub":               Some(forgedomain.ForgeTypeGitHub),
			"gitlab":               Some(forgedomain.ForgeTypeGitLab),
			"GitLab":               Some(forgedomain.ForgeTypeGitLab),
			"gitea":                Some(forgedomain.ForgeTypeGitea),
			"Gitea":                Some(forgedomain.ForgeTypeGitea),
		}
		for give, want := range tests {
			have, err := forgedomain.ParseForgeType(give)
			must.NoError(t, err)
			must.Eq(t, want, have)
		}
	})

	t.Run("invalid content", func(t *testing.T) {
		t.Parallel()
		_, err := forgedomain.ParseForgeType("zonk")
		must.Error(t, err)
	})
}
