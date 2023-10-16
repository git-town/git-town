package bitbucket_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/giturl"
	"github.com/git-town/git-town/v9/src/hosting/bitbucket"
	"github.com/git-town/git-town/v9/src/hosting/common"
	"github.com/shoenig/test/must"
)

// emptySHAForBranch is a dummy implementation for hosting.SHAForBranchfunc to be used in tests.
func emptySHAForBranch(domain.BranchName) (domain.SHA, error) {
	return domain.EmptySHA(), nil
}

func TestNewBitbucketConnector(t *testing.T) {
	t.Parallel()

	t.Run("Bitbucket SaaS", func(t *testing.T) {
		t.Parallel()
		have, err := bitbucket.NewConnector(bitbucket.NewConnectorArgs{
			HostingService:  config.HostingNone,
			OriginURL:       giturl.Parse("username@bitbucket.org:git-town/docs.git"),
			GetSHAForBranch: emptySHAForBranch,
		})
		must.NoError(t, err)
		wantConfig := common.Config{
			APIToken:     "",
			Hostname:     "bitbucket.org",
			Organization: "git-town",
			Repository:   "docs",
		}
		must.EqOp(t, wantConfig, have.Config)
	})

	t.Run("hosted service type provided manually", func(t *testing.T) {
		t.Parallel()
		have, err := bitbucket.NewConnector(bitbucket.NewConnectorArgs{
			HostingService:  config.HostingBitbucket,
			OriginURL:       giturl.Parse("git@custom-url.com:git-town/docs.git"),
			GetSHAForBranch: emptySHAForBranch,
		})
		must.NoError(t, err)
		wantConfig := common.Config{
			APIToken:     "",
			Hostname:     "custom-url.com",
			Organization: "git-town",
			Repository:   "docs",
		}
		must.EqOp(t, wantConfig, have.Config)
	})

	t.Run("repo is hosted by another hosting service --> no connector", func(t *testing.T) {
		t.Parallel()
		have, err := bitbucket.NewConnector(bitbucket.NewConnectorArgs{
			HostingService:  config.HostingNone,
			OriginURL:       giturl.Parse("git@github.com:git-town/git-town.git"),
			GetSHAForBranch: emptySHAForBranch,
		})
		must.Nil(t, have)
		must.NoError(t, err)
	})

	t.Run("no origin remote --> no connector", func(t *testing.T) {
		t.Parallel()
		var originURL *giturl.Parts
		have, err := bitbucket.NewConnector(bitbucket.NewConnectorArgs{
			HostingService:  config.HostingNone,
			OriginURL:       originURL,
			GetSHAForBranch: emptySHAForBranch,
		})
		must.Nil(t, have)
		must.NoError(t, err)
	})
}
