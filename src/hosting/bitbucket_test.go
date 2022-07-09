package hosting_test

import (
	"testing"

	"github.com/git-town/git-town/v7/src/hosting"
	"github.com/stretchr/testify/assert"
)

//nolint:paralleltest  // mocks HTTP
func TestLoadBitbucket(t *testing.T) {
	driver := hosting.NewBitbucketDriver(mockConfig{
		hostingService: "bitbucket",
		originURL:      "git@self-hosted-bitbucket.com:git-town/git-town.git",
	}, nil)
	assert.NotNil(t, driver)
	assert.Equal(t, "Bitbucket", driver.HostingServiceName())
	assert.Equal(t, "https://self-hosted-bitbucket.com/git-town/git-town", driver.RepositoryURL())
}

//nolint:paralleltest  // mocks HTTP
func TestLoadBitbucket_customHostName(t *testing.T) {
	driver := hosting.NewBitbucketDriver(mockConfig{
		originURL:      "git@my-ssh-identity.com:git-town/git-town.git",
		originOverride: "bitbucket.org",
	}, nil)
	assert.NotNil(t, driver)
	assert.Equal(t, "Bitbucket", driver.HostingServiceName())
	assert.Equal(t, "https://bitbucket.org/git-town/git-town", driver.RepositoryURL())
}

//nolint:paralleltest  // mocks HTTP
func TestLoadBitbucket_customUserName(t *testing.T) {
	driver := hosting.NewBitbucketDriver(mockConfig{
		hostingService: "bitbucket",
		originURL:      "username@bitbucket.org:git-town/git-town.git",
	}, nil)
	assert.NotNil(t, driver)
	assert.Equal(t, "Bitbucket", driver.HostingServiceName())
	assert.Equal(t, "https://bitbucket.org/git-town/git-town", driver.RepositoryURL())
}
