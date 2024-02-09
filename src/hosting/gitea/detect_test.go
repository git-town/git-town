package gitea_test

import (
	"testing"

	"github.com/git-town/git-town/v12/src/config/configdomain"
	"github.com/git-town/git-town/v12/src/git/giturl"
	"github.com/git-town/git-town/v12/src/hosting/gitea"
	"github.com/shoenig/test/must"
)

func TestDetect(t *testing.T) {
	t.Parallel()

	t.Run("hosted service type provided manually", func(t *testing.T) {
		t.Parallel()
		must.True(t, gitea.Detect(giturl.Parse("git@custom-url.com:git-town/docs.git"), configdomain.HostingPlatformGitea))
	})

	t.Run("repo is hosted by another hosting platform", func(t *testing.T) {
		t.Parallel()
		must.False(t, gitea.Detect(giturl.Parse("git@github.com:git-town/git-town.git"), configdomain.HostingPlatformNone))
	})

	t.Run("no origin remote", func(t *testing.T) {
		t.Parallel()
		var originURL *giturl.Parts
		must.False(t, gitea.Detect(originURL, configdomain.HostingPlatformNone))
	})
}
