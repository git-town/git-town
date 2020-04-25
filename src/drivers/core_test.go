package drivers_test

import (
	"testing"

	"github.com/git-town/git-town/src/drivers"
	"github.com/stretchr/testify/assert"
)

func TestGetDriver_DriverType_Bitbucket(t *testing.T) {
	driver := drivers.GetDriver(drivers.DriverOptions{
		DriverType: "bitbucket",
		OriginURL:  "git@self-hosted-bitbucket.com:git-town/git-town.git",
	})

	assert.NotNil(t, driver)
	assert.Equal(t, "Bitbucket", driver.HostingServiceName())
	assert.Equal(t, "https://self-hosted-bitbucket.com/git-town/git-town", driver.GetRepositoryURL())
}

func TestGetDriver_DriverType_GitHub(t *testing.T) {
	driver := drivers.GetDriver(drivers.DriverOptions{
		DriverType: "github",
		OriginURL:  "git@self-hosted-github.com:git-town/git-town.git",
	})

	assert.NotNil(t, driver)
	assert.Equal(t, "GitHub", driver.HostingServiceName())
	assert.Equal(t, "https://self-hosted-github.com/git-town/git-town", driver.GetRepositoryURL())
}

func TestGetDriver_DriverType_GitLab(t *testing.T) {
	driver := drivers.GetDriver(drivers.DriverOptions{
		DriverType: "gitlab",
		OriginURL:  "git@self-hosted-gitlab.com:git-town/git-town.git",
	})

	assert.NotNil(t, driver)
	assert.Equal(t, "GitLab", driver.HostingServiceName())
	assert.Equal(t, "https://self-hosted-gitlab.com/git-town/git-town", driver.GetRepositoryURL())
}

func TestGetDriver_OriginHostname_Bitbucket(t *testing.T) {
	driver := drivers.GetDriver(drivers.DriverOptions{
		OriginURL:      "git@my-ssh-identity.com:git-town/git-town.git",
		OriginHostname: "bitbucket.org",
	})

	assert.NotNil(t, driver)
	assert.Equal(t, "Bitbucket", driver.HostingServiceName())
	assert.Equal(t, "https://bitbucket.org/git-town/git-town", driver.GetRepositoryURL())
}

func TestGetDriver_OriginHostname_GitHub(t *testing.T) {
	driver := drivers.GetDriver(drivers.DriverOptions{
		OriginURL:      "git@my-ssh-identity.com:git-town/git-town.git",
		OriginHostname: "github.com",
	})

	assert.NotNil(t, driver)
	assert.Equal(t, "GitHub", driver.HostingServiceName())
	assert.Equal(t, "https://github.com/git-town/git-town", driver.GetRepositoryURL())
}

func TestGetDriver_OriginHostname_GitLab(t *testing.T) {
	driver := drivers.GetDriver(drivers.DriverOptions{
		OriginURL:      "git@my-ssh-identity.com:git-town/git-town.git",
		OriginHostname: "gitlab.com",
	})

	assert.NotNil(t, driver)
	assert.Equal(t, "GitLab", driver.HostingServiceName())
	assert.Equal(t, "https://gitlab.com/git-town/git-town", driver.GetRepositoryURL())
}
