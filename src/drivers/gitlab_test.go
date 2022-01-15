package drivers_test

import (
	"testing"

	"github.com/git-town/git-town/v7/src/drivers"
	"github.com/stretchr/testify/assert"
)

//nolint:paralleltest
func TestLoadGitLab(t *testing.T) {
	driver := drivers.LoadGitlab(mockConfig{
		codeHostingDriverName: "gitlab",
		remoteOriginURL:       "git@self-hosted-gitlab.com:git-town/git-town.git",
	})
	assert.NotNil(t, driver)
	assert.Equal(t, "GitLab", driver.HostingServiceName())
	assert.Equal(t, "https://self-hosted-gitlab.com/git-town/git-town", driver.RepositoryURL())
}

//nolint:paralleltest
func TestLoadGitLab_customHostName(t *testing.T) {
	driver := drivers.LoadGitlab(mockConfig{
		remoteOriginURL:    "git@my-ssh-identity.com:git-town/git-town.git",
		configuredHostName: "gitlab.com",
	})
	assert.NotNil(t, driver)
	assert.Equal(t, "GitLab", driver.HostingServiceName())
	assert.Equal(t, "https://gitlab.com/git-town/git-town", driver.RepositoryURL())
}
