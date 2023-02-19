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
		bitbucketDriver := hosting.NewBitbucketDriver(*url, repoConfig, nil)
		assert.NotNil(t, bitbucketDriver)
		assert.Equal(t, "Bitbucket", bitbucketDriver.HostingServiceName())
		assert.Equal(t, "https://self-hosted-bitbucket.com/git-town/git-town", bitbucketDriver.RepositoryURL())
	})

	t.Run("custom hostname", func(t *testing.T) {
		t.Parallel()
		config := mockRepoConfig{
			originURL:      "git@my-ssh-identity.com:git-town/git-town.git",
			originOverride: "bitbucket.org",
		}
		url := giturl.Parse(config.originURL)
		driver := hosting.NewBitbucketDriver(*url, config, nil)
		assert.NotNil(t, driver)
		assert.Equal(t, "Bitbucket", driver.HostingServiceName())
		assert.Equal(t, "https://bitbucket.org/git-town/git-town", driver.RepositoryURL())
	})

	t.Run("custom username", func(t *testing.T) {
		t.Parallel()
		config := mockRepoConfig{
			hostingService: "bitbucket",
			originURL:      "username@bitbucket.org:git-town/git-town.git",
		}
		url := giturl.Parse(config.originURL)
		driver := hosting.NewBitbucketDriver(*url, config, nil)
		assert.NotNil(t, driver)
		assert.Equal(t, "Bitbucket", driver.HostingServiceName())
		assert.Equal(t, "https://bitbucket.org/git-town/git-town", driver.RepositoryURL())
	})
}
