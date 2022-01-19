package drivers_test

import (
	"testing"

	"github.com/git-town/git-town/v7/src/drivers"
	"github.com/stretchr/testify/assert"
)

//nolint:paralleltest  // mocks HTTP
func TestLoadBitbucket(t *testing.T) {
	driver := drivers.LoadBitbucket(mockConfig{
		codeHostingDriverName: "bitbucket",
		remoteOriginURL:       "git@self-hosted-bitbucket.com:git-town/git-town.git",
	}, nil)
	assert.NotNil(t, driver)
	assert.Equal(t, "Bitbucket", driver.HostingServiceName())
	assert.Equal(t, "https://self-hosted-bitbucket.com/git-town/git-town", driver.RepositoryURL())
}

//nolint:paralleltest  // mocks HTTP
func TestLoadBitbucket_customHostName(t *testing.T) {
	driver := drivers.LoadBitbucket(mockConfig{
		remoteOriginURL: "git@my-ssh-identity.com:git-town/git-town.git",
		manualHostName:  "bitbucket.org",
	}, nil)
	assert.NotNil(t, driver)
	assert.Equal(t, "Bitbucket", driver.HostingServiceName())
	assert.Equal(t, "https://bitbucket.org/git-town/git-town", driver.RepositoryURL())
}
