package hosting_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/giturl"
	"github.com/git-town/git-town/v9/src/hosting"
	"github.com/stretchr/testify/assert"
)

func TestNewBitbucketConnector(t *testing.T) {
	t.Parallel()
	t.Run("Bitbucket SaaS", func(t *testing.T) {
		t.Parallel()
		have, err := hosting.NewBitbucketConnector(hosting.NewBitbucketConnectorArgs{
			HostingService:  config.HostingNone,
			OriginURL:       giturl.Parse("username@bitbucket.org:git-town/docs.git"),
			GetSHAForBranch: emptySHAForBranch,
		})
		assert.NoError(t, err)
		wantConfig := hosting.CommonConfig{
			APIToken:     "",
			Hostname:     "bitbucket.org",
			Organization: "git-town",
			Repository:   "docs",
		}
		assert.Equal(t, wantConfig, have.CommonConfig)
	})
	t.Run("hosted service type provided manually", func(t *testing.T) {
		t.Parallel()
		have, err := hosting.NewBitbucketConnector(hosting.NewBitbucketConnectorArgs{
			HostingService:  config.HostingBitbucket,
			OriginURL:       giturl.Parse("git@custom-url.com:git-town/docs.git"),
			GetSHAForBranch: emptySHAForBranch,
		})
		assert.NoError(t, err)
		wantConfig := hosting.CommonConfig{
			APIToken:     "",
			Hostname:     "custom-url.com",
			Organization: "git-town",
			Repository:   "docs",
		}
		assert.Equal(t, wantConfig, have.CommonConfig)
	})
	t.Run("repo is hosted by another hosting service --> no connector", func(t *testing.T) {
		t.Parallel()
		have, err := hosting.NewBitbucketConnector(hosting.NewBitbucketConnectorArgs{
			HostingService:  config.HostingNone,
			OriginURL:       giturl.Parse("git@github.com:git-town/git-town.git"),
			GetSHAForBranch: emptySHAForBranch,
		})
		assert.Nil(t, have)
		assert.NoError(t, err)
	})
	t.Run("no origin remote --> no connector", func(t *testing.T) {
		t.Parallel()
		var originURL *giturl.Parts
		have, err := hosting.NewBitbucketConnector(hosting.NewBitbucketConnectorArgs{
			HostingService:  config.HostingNone,
			OriginURL:       originURL,
			GetSHAForBranch: emptySHAForBranch,
		})
		assert.Nil(t, have)
		assert.NoError(t, err)
	})
}
