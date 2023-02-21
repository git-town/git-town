package hosting_test

import (
	"testing"

	"github.com/git-town/git-town/v7/src/giturl"
	"github.com/git-town/git-town/v7/src/hosting"
	"github.com/stretchr/testify/assert"
)

func TestNewBitbucketDriver(t *testing.T) {
	t.Parallel()
	t.Run("normal example", func(t *testing.T) {
		t.Parallel()
		repoConfig := mockRepoConfig{
			hostingService: "bitbucket",
			originURL:      "git@self-hosted-bitbucket.com:git-town/git-town.git",
		}
		url := giturl.Parse(repoConfig.originURL)
		connector := hosting.NewBitbucketConnector(*url, repoConfig, nil)
		assert.NotNil(t, connector)
		assert.Equal(t, "Bitbucket", connector.HostingServiceName())
		assert.Equal(t, "https://self-hosted-bitbucket.com/git-town/git-town", connector.RepositoryURL())
	})

	t.Run("custom hostname", func(t *testing.T) {
		t.Parallel()
		repoConfig := mockRepoConfig{
			originURL:      "git@my-ssh-identity.com:git-town/git-town.git",
			originOverride: "bitbucket.org",
		}
		url := giturl.Parse(repoConfig.originURL)
		connector := hosting.NewBitbucketConnector(*url, repoConfig, nil)
		assert.NotNil(t, connector)
		assert.Equal(t, "Bitbucket", connector.HostingServiceName())
		assert.Equal(t, "https://bitbucket.org/git-town/git-town", connector.RepositoryURL())
	})

	t.Run("custom username", func(t *testing.T) {
		t.Parallel()
		repoConfig := mockRepoConfig{
			hostingService: "bitbucket",
			originURL:      "username@bitbucket.org:git-town/git-town.git",
		}
		url := giturl.Parse(repoConfig.originURL)
		connector := hosting.NewBitbucketConnector(*url, repoConfig, nil)
		assert.NotNil(t, connector)
		assert.Equal(t, "Bitbucket", connector.HostingServiceName())
		assert.Equal(t, "https://bitbucket.org/git-town/git-town", connector.RepositoryURL())
	})
}
