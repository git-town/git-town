package hosting_test

import (
	"testing"

	"github.com/git-town/git-town/v7/src/hosting"
	"github.com/stretchr/testify/assert"
)

//nolint:paralleltest  // mocks HTTP
func TestLoadGitLab(t *testing.T) {
	driver := hosting.NewGitlabDriver(mockConfig{
		hostingService: "gitlab",
		originURL:      "git@self-hosted-gitlab.com:git-town/git-town.git",
	})
	assert.NotNil(t, driver)
	assert.Equal(t, "GitLab", driver.HostingServiceName())
	assert.Equal(t, "https://self-hosted-gitlab.com/git-town/git-town", driver.RepositoryURL())
}

//nolint:paralleltest  // mocks HTTP
func TestLoadGitLab_customHostName(t *testing.T) {
	driver := hosting.NewGitlabDriver(mockConfig{
		originURL:      "git@my-ssh-identity.com:git-town/git-town.git",
		originOverride: "gitlab.com",
	})
	assert.NotNil(t, driver)
	assert.Equal(t, "GitLab", driver.HostingServiceName())
	assert.Equal(t, "https://gitlab.com/git-town/git-town", driver.RepositoryURL())
}
