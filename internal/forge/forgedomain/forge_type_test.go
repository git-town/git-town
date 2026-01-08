package forgedomain_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestParseForgeType(t *testing.T) {
	t.Parallel()

	t.Run("acceptable content", func(t *testing.T) {
		t.Parallel()
		tests := map[string]Option[forgedomain.ForgeType]{
			"":                     None[forgedomain.ForgeType](),
			"azuredevops":          Some(forgedomain.ForgeTypeAzuredevops),
			"AzureDevOps":          Some(forgedomain.ForgeTypeAzuredevops),
			"AZUREDEVOPS":          Some(forgedomain.ForgeTypeAzuredevops),
			"bitbucket":            Some(forgedomain.ForgeTypeBitbucket),
			"BitBucket":            Some(forgedomain.ForgeTypeBitbucket),
			"BITBUCKET":            Some(forgedomain.ForgeTypeBitbucket),
			"bitbucket-datacenter": Some(forgedomain.ForgeTypeBitbucketDatacenter),
			"BitBucket-Datacenter": Some(forgedomain.ForgeTypeBitbucketDatacenter),
			"BITBUCKET-DATACENTER": Some(forgedomain.ForgeTypeBitbucketDatacenter),
			"forgejo":              Some(forgedomain.ForgeTypeForgejo),
			"Forgejo":              Some(forgedomain.ForgeTypeForgejo),
			"ForgeJo":              Some(forgedomain.ForgeTypeForgejo),
			"FORGEJO":              Some(forgedomain.ForgeTypeForgejo),
			"github":               Some(forgedomain.ForgeTypeGithub),
			"GitHub":               Some(forgedomain.ForgeTypeGithub),
			"gitlab":               Some(forgedomain.ForgeTypeGitlab),
			"GitLab":               Some(forgedomain.ForgeTypeGitlab),
			"gitea":                Some(forgedomain.ForgeTypeGitea),
			"Gitea":                Some(forgedomain.ForgeTypeGitea),
		}
		for give, want := range tests {
			have, err := forgedomain.ParseForgeType(give, "test")
			must.NoError(t, err)
			must.Eq(t, want, have)
		}
	})

	t.Run("invalid content", func(t *testing.T) {
		t.Parallel()
		_, err := forgedomain.ParseForgeType("zonk", "test")
		must.Error(t, err)
	})
}
