package drivers_test

import (
	"testing"

	"github.com/git-town/git-town/src/drivers"
	"github.com/stretchr/testify/assert"
)

type mockGitlabConfig struct {
	codeHostingDriverName string
	remoteOriginURL       string
	configuredHostName    string
}

func (mgc mockGitlabConfig) GetCodeHostingDriverName() string {
	return mgc.codeHostingDriverName
}
func (mgc mockGitlabConfig) GetRemoteOriginURL() string {
	return mgc.remoteOriginURL
}
func (mgc mockGitlabConfig) GetCodeHostingOriginHostname() string {
	return mgc.configuredHostName
}

func TestGetDriver_DriverType_GitLab(t *testing.T) {
	driver := drivers.LoadGitlab(mockGitlabConfig{
		codeHostingDriverName: "gitlab",
		remoteOriginURL:       "git@self-hosted-gitlab.com:git-town/git-town.git",
	})
	assert.NotNil(t, driver)
	assert.Equal(t, "GitLab", driver.HostingServiceName())
	assert.Equal(t, "https://self-hosted-gitlab.com/git-town/git-town", driver.GetRepositoryURL())
}

func TestGetDriver_OriginHostname_GitLab(t *testing.T) {
	driver := drivers.LoadGitlab(mockGitlabConfig{
		remoteOriginURL:    "git@my-ssh-identity.com:git-town/git-town.git",
		configuredHostName: "gitlab.com",
	})
	assert.NotNil(t, driver)
	assert.Equal(t, "GitLab", driver.HostingServiceName())
	assert.Equal(t, "https://gitlab.com/git-town/git-town", driver.GetRepositoryURL())
}
