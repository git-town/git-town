package drivers_test

import (
	"testing"

	"github.com/Originate/git-town/src/drivers"
	"github.com/stretchr/testify/assert"
)

func TestGetDriver_DriverOverride_Bitbucket(t *testing.T) {
	result := drivers.GetDriver(drivers.DriverOptions{
		DriverType: "bitbucket",
		OriginURL:  "git@self-hosted-bitbucket.com:Originate/git-town.git",
	})
	assert.NotNil(t, result)
	assert.Equal(t, "Bitbucket", result.HostingServiceName())
	assert.Equal(t, "https://self-hosted-bitbucket.com/Originate/git-town", result.GetRepositoryURL())
}

func TestGetDriver_DriverOverride_GitHub(t *testing.T) {
	result := drivers.GetDriver(drivers.DriverOptions{
		DriverType: "github",
		OriginURL:  "git@self-hosted-github.com:Originate/git-town.git",
	})
	assert.NotNil(t, result)
	assert.Equal(t, "GitHub", result.HostingServiceName())
	assert.Equal(t, "https://self-hosted-github.com/Originate/git-town", result.GetRepositoryURL())
}

func TestGetDriver_DriverOverride_GitLab(t *testing.T) {
	result := drivers.GetDriver(drivers.DriverOptions{
		DriverType: "gitlab",
		OriginURL:  "git@self-hosted-gitlab.com:Originate/git-town.git",
	})
	assert.NotNil(t, result)
	assert.Equal(t, "GitLab", result.HostingServiceName())
	assert.Equal(t, "https://self-hosted-gitlab.com/Originate/git-town", result.GetRepositoryURL())
}

func TestGetDriver_OriginHostnameOverride_Bitbucket(t *testing.T) {
	result := drivers.GetDriver(drivers.DriverOptions{
		OriginURL:      "git@my-ssh-identity.com:Originate/git-town.git",
		OriginHostname: "bitbucket.org",
	})
	assert.NotNil(t, result)
	assert.Equal(t, "Bitbucket", result.HostingServiceName())
	assert.Equal(t, "https://bitbucket.org/Originate/git-town", result.GetRepositoryURL())
}

func TestGetDriver_OriginHostnameOverride_GitHub(t *testing.T) {
	result := drivers.GetDriver(drivers.DriverOptions{
		OriginURL:      "git@my-ssh-identity.com:Originate/git-town.git",
		OriginHostname: "github.com",
	})
	assert.NotNil(t, result)
	assert.Equal(t, "GitHub", result.HostingServiceName())
	assert.Equal(t, "https://github.com/Originate/git-town", result.GetRepositoryURL())
}

func TestGetDriver_OriginHostnameOverride_GitLab(t *testing.T) {
	result := drivers.GetDriver(drivers.DriverOptions{
		OriginURL:      "git@my-ssh-identity.com:Originate/git-town.git",
		OriginHostname: "gitlab.com",
	})
	assert.NotNil(t, result)
	assert.Equal(t, "GitLab", result.HostingServiceName())
	assert.Equal(t, "https://gitlab.com/Originate/git-town", result.GetRepositoryURL())
}
