package bitbucket_test

import (
	"testing"

	"github.com/git-town/git-town/v13/src/config/configdomain"
	"github.com/git-town/git-town/v13/src/git/giturl"
	"github.com/git-town/git-town/v13/src/hosting/bitbucket"
	"github.com/shoenig/test/must"
)

func TestDetect(t *testing.T) {
	t.Parallel()

	t.Run("Bitbucket SaaS", func(t *testing.T) {
		t.Parallel()
		t.Run("Bitbucket SaaS", func(t *testing.T) {
			t.Parallel()
			must.True(t, bitbucket.Detect(giturl.Parse("username@bitbucket.org:git-town/docs.git"), configdomain.HostingPlatformNone))
		})
		t.Run("hosted service type provided manually", func(t *testing.T) {
			t.Parallel()
			must.True(t, bitbucket.Detect(giturl.Parse("git@custom-url.com:git-town/docs.git"), configdomain.HostingPlatformBitbucket))
		})
		t.Run("repo is hosted by another hosting platform", func(t *testing.T) {
			t.Parallel()
			must.False(t, bitbucket.Detect(giturl.Parse("git@github.com:git-town/git-town.git"), configdomain.HostingPlatformNone))
		})
		t.Run("no origin remote", func(t *testing.T) {
			t.Parallel()
			var originURL *giturl.Parts
			must.False(t, bitbucket.Detect(originURL, configdomain.HostingPlatformNone))
		})
	})
}
