package drivers_test

import (
	"testing"

	"github.com/git-town/git-town/src/drivers"
	"github.com/stretchr/testify/assert"
)

type mockBitbucketConfig struct {
	codeHostingDriverName string
	remoteOriginURL       string
	configuredHostName    string
}

func (mbc mockBitbucketConfig) GetCodeHostingDriverName() string {
	return mbc.codeHostingDriverName
}
func (mbc mockBitbucketConfig) GetRemoteOriginURL() string {
	return mbc.remoteOriginURL
}

func (mbc mockBitbucketConfig) GetCodeHostingOriginHostname() string {
	return mbc.configuredHostName
}

func TestTryUseBitbucket(t *testing.T) {
	driver := drivers.TryUseBitbucket(mockBitbucketConfig{
		codeHostingDriverName: "bitbucket",
		remoteOriginURL:       "git@self-hosted-bitbucket.com:git-town/git-town.git",
	})
	assert.NotNil(t, driver)
	assert.Equal(t, "Bitbucket", driver.HostingServiceName())
	assert.Equal(t, "https://self-hosted-bitbucket.com/git-town/git-town", driver.GetRepositoryURL())
}

func TestTryUseBitbucket_customHostName(t *testing.T) {
	driver := drivers.TryUseBitbucket(mockBitbucketConfig{
		remoteOriginURL:    "git@my-ssh-identity.com:git-town/git-town.git",
		configuredHostName: "bitbucket.org",
	})
	assert.NotNil(t, driver)
	assert.Equal(t, "Bitbucket", driver.HostingServiceName())
	assert.Equal(t, "https://bitbucket.org/git-town/git-town", driver.GetRepositoryURL())
}
