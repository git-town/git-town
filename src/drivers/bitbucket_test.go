package drivers_test

import (
	"testing"

	"github.com/git-town/git-town/src/drivers"
	"github.com/stretchr/testify/assert"
)

func TestLoadBitbucket(t *testing.T) {
	driver := drivers.LoadBitbucket(mockConfig{
		codeHostingDriverName: "bitbucket",
		remoteOriginURL:       "git@self-hosted-bitbucket.com:git-town/git-town.git",
	}, nil)
	assert.NotNil(t, driver)
	assert.Equal(t, "Bitbucket", driver.HostingServiceName())
	assert.Equal(t, "https://self-hosted-bitbucket.com/git-town/git-town", driver.RepositoryURL())
}

func TestLoadBitbucket_customHostName(t *testing.T) {
	driver := drivers.LoadBitbucket(mockConfig{
		remoteOriginURL:    "git@my-ssh-identity.com:git-town/git-town.git",
		configuredHostName: "bitbucket.org",
	}, nil)
	assert.NotNil(t, driver)
	assert.Equal(t, "Bitbucket", driver.HostingServiceName())
	assert.Equal(t, "https://bitbucket.org/git-town/git-town", driver.RepositoryURL())
}
