package github_test

import (
	"testing"

	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/git/giturl"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/git-town/git-town/v14/src/hosting/github"
	"github.com/shoenig/test/must"
)

func TestDetect(t *testing.T) {
	t.Parallel()

	t.Run("GitHub SaaS", func(t *testing.T) {
		t.Parallel()
		must.True(t, github.Detect(giturl.Parse("git@github.com:git-town/docs.git"), None[configdomain.HostingPlatform]()))
	})

	t.Run("hosted service type provided manually", func(t *testing.T) {
		t.Parallel()
		must.True(t, github.Detect(giturl.Parse("git@custom-url.com:git-town/docs.git"), Some(configdomain.HostingPlatformGitHub)))
	})

	t.Run("repo is hosted by another hosting platform", func(t *testing.T) {
		t.Parallel()
		must.False(t, github.Detect(giturl.Parse("git@gitlab.com:git-town/git-town.git"), None[configdomain.HostingPlatform]()))
	})

	t.Run("no origin remote", func(t *testing.T) {
		t.Parallel()
		var originURL *giturl.Parts
		must.False(t, github.Detect(originURL, None[configdomain.HostingPlatform]()))
	})
}
