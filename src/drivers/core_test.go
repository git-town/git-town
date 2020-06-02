package drivers_test

import (
	"testing"

	gtmocks "github.com/git-town/git-town/mocks"
	"github.com/git-town/git-town/src/drivers"
	"github.com/stretchr/testify/assert"
)

func TestGetDriver_DriverType_Bitbucket(t *testing.T) {
	var mockedGitConfig = new(gtmocks.ConfigurationInterface)
	mockedGitConfig.On("GetCodeHostingOriginHostname").Return("")
	mockedGitConfig.On("GetCodeHostingDriverName").Return("bitbucket")
	mockedGitConfig.On("GetRemoteOriginURL").Return("git@self-hosted-bitbucket.com:git-town/git-town.git")
	drivers.GitConfig = mockedGitConfig
	driver := drivers.GetDriver()
	assert.NotNil(t, driver)
	assert.Equal(t, "Bitbucket", driver.HostingServiceName())
	assert.Equal(t, "https://self-hosted-bitbucket.com/git-town/git-town", driver.GetRepositoryURL())
	mockedGitConfig.AssertExpectations(t)
}

func TestGetDriver_DriverType_GitHub(t *testing.T) {
	var mockedGitConfig = new(gtmocks.ConfigurationInterface)
	mockedGitConfig.On("GetCodeHostingOriginHostname").Return("")
	mockedGitConfig.On("GetCodeHostingDriverName").Return("github")
	mockedGitConfig.On("GetGitHubToken").Return("")
	mockedGitConfig.On("GetRemoteOriginURL").Return("git@self-hosted-github.com:git-town/git-town.git")
	drivers.GitConfig = mockedGitConfig
	driver := drivers.GetDriver()
	assert.NotNil(t, driver)
	assert.Equal(t, "GitHub", driver.HostingServiceName())
	assert.Equal(t, "https://self-hosted-github.com/git-town/git-town", driver.GetRepositoryURL())
	mockedGitConfig.AssertExpectations(t)
}

func TestGetDriver_DriverType_Gitea(t *testing.T) {
	var mockedGitConfig = new(gtmocks.ConfigurationInterface)
	mockedGitConfig.On("GetCodeHostingOriginHostname").Return("")
	mockedGitConfig.On("GetCodeHostingDriverName").Return("gitea")
	mockedGitConfig.On("GetGiteaToken").Return("")
	mockedGitConfig.On("GetRemoteOriginURL").Return("git@self-hosted-gitea.com:git-town/git-town.git")
	drivers.GitConfig = mockedGitConfig
	driver := drivers.GetDriver()
	assert.NotNil(t, driver)
	assert.Equal(t, "Gitea", driver.HostingServiceName())
	assert.Equal(t, "https://self-hosted-gitea.com/git-town/git-town", driver.GetRepositoryURL())
	mockedGitConfig.AssertExpectations(t)
}

func TestGetDriver_DriverType_GitLab(t *testing.T) {
	var mockedGitConfig = new(gtmocks.ConfigurationInterface)
	mockedGitConfig.On("GetCodeHostingOriginHostname").Return("")
	mockedGitConfig.On("GetCodeHostingDriverName").Return("gitlab")
	mockedGitConfig.On("GetRemoteOriginURL").Return("git@self-hosted-gitlab.com:git-town/git-town.git")
	drivers.GitConfig = mockedGitConfig
	driver := drivers.GetDriver()
	assert.NotNil(t, driver)
	assert.Equal(t, "GitLab", driver.HostingServiceName())
	assert.Equal(t, "https://self-hosted-gitlab.com/git-town/git-town", driver.GetRepositoryURL())
	mockedGitConfig.AssertExpectations(t)
}

func TestGetDriver_OriginHostname_Bitbucket(t *testing.T) {
	var mockedGitConfig = new(gtmocks.ConfigurationInterface)
	mockedGitConfig.On("GetCodeHostingOriginHostname").Return("bitbucket.org")
	mockedGitConfig.On("GetCodeHostingDriverName").Return("")
	mockedGitConfig.On("GetRemoteOriginURL").Return("git@my-ssh-identity.com:git-town/git-town.git")
	drivers.GitConfig = mockedGitConfig
	driver := drivers.GetDriver()
	assert.NotNil(t, driver)
	assert.Equal(t, "Bitbucket", driver.HostingServiceName())
	assert.Equal(t, "https://bitbucket.org/git-town/git-town", driver.GetRepositoryURL())
	mockedGitConfig.AssertExpectations(t)
}

func TestGetDriver_OriginHostname_GitHub(t *testing.T) {
	var mockedGitConfig = new(gtmocks.ConfigurationInterface)
	mockedGitConfig.On("GetCodeHostingOriginHostname").Return("github.com")
	mockedGitConfig.On("GetCodeHostingDriverName").Return("")
	mockedGitConfig.On("GetGitHubToken").Return("")
	mockedGitConfig.On("GetRemoteOriginURL").Return("git@my-ssh-identity.com:git-town/git-town.git")
	drivers.GitConfig = mockedGitConfig
	driver := drivers.GetDriver()
	assert.NotNil(t, driver)
	assert.Equal(t, "GitHub", driver.HostingServiceName())
	assert.Equal(t, "https://github.com/git-town/git-town", driver.GetRepositoryURL())
	mockedGitConfig.AssertExpectations(t)
}

func TestGetDriver_OriginHostname_Gitea(t *testing.T) {
	var mockedGitConfig = new(gtmocks.ConfigurationInterface)
	mockedGitConfig.On("GetCodeHostingOriginHostname").Return("gitea.com")
	mockedGitConfig.On("GetCodeHostingDriverName").Return("")
	mockedGitConfig.On("GetGiteaToken").Return("")
	mockedGitConfig.On("GetRemoteOriginURL").Return("git@my-ssh-identity.com:git-town/git-town.git")
	drivers.GitConfig = mockedGitConfig
	driver := drivers.GetDriver()
	assert.NotNil(t, driver)
	assert.Equal(t, "Gitea", driver.HostingServiceName())
	assert.Equal(t, "https://gitea.com/git-town/git-town", driver.GetRepositoryURL())
	mockedGitConfig.AssertExpectations(t)
}

func TestGetDriver_OriginHostname_GitLab(t *testing.T) {
	var mockedGitConfig = new(gtmocks.ConfigurationInterface)
	mockedGitConfig.On("GetCodeHostingOriginHostname").Return("gitlab.com")
	mockedGitConfig.On("GetCodeHostingDriverName").Return("")
	mockedGitConfig.On("GetRemoteOriginURL").Return("git@my-ssh-identity.com:git-town/git-town.git")
	drivers.GitConfig = mockedGitConfig
	driver := drivers.GetDriver()
	assert.NotNil(t, driver)
	assert.Equal(t, "GitLab", driver.HostingServiceName())
	assert.Equal(t, "https://gitlab.com/git-town/git-town", driver.GetRepositoryURL())
	mockedGitConfig.AssertExpectations(t)
}
