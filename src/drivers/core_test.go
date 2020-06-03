package drivers_test

import (
	"testing"

	"github.com/git-town/git-town/src/drivers"
	"github.com/stretchr/testify/assert"
)

func TestGetDriver_DriverType_GitLab(t *testing.T) {
	driver := drivers.GetDriver(drivers.DriverOptions{
		DriverType: "gitlab",
		OriginURL:  "git@self-hosted-gitlab.com:git-town/git-town.git",
	})
	assert.NotNil(t, driver)
	assert.Equal(t, "GitLab", driver.HostingServiceName())
	assert.Equal(t, "https://self-hosted-gitlab.com/git-town/git-town", driver.GetRepositoryURL())
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
