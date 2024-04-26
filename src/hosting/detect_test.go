package hosting_test

import (
	"testing"

	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/git/giturl"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/git-town/git-town/v14/src/hosting"
	"github.com/shoenig/test/must"
)

func TestDetect(t *testing.T) {
	t.Parallel()
	t.Run("BitBucket SAAS, no override", func(t *testing.T) {
		t.Parallel()
		url := giturl.Parse("username@bitbucket.org:git-town/docs.git")
		have := hosting.Detect(url, None[configdomain.HostingPlatform]())
		want := Some(configdomain.HostingPlatformBitbucket)
		must.Eq(t, want, have)
	})
	t.Run("custom URL, override to BitBucket", func(t *testing.T) {
		t.Parallel()
		url := giturl.Parse("username@custom.org:git-town/docs.git")
		have := hosting.Detect(url, Some(configdomain.HostingPlatformBitbucket))
		want := Some(configdomain.HostingPlatformBitbucket)
		must.Eq(t, want, have)
	})
	t.Run("Gitea SAAS, no override", func(t *testing.T) {
		t.Parallel()
		url := giturl.Parse("username@gitea.com:git-town/docs.git")
		have := hosting.Detect(url, None[configdomain.HostingPlatform]())
		want := Some(configdomain.HostingPlatformGitea)
		must.Eq(t, want, have)
	})
	t.Run("custom URL, override to Gitea", func(t *testing.T) {
		t.Parallel()
		url := giturl.Parse("username@custom.org:git-town/docs.git")
		have := hosting.Detect(url, Some(configdomain.HostingPlatformGitea))
		want := Some(configdomain.HostingPlatformGitea)
		must.Eq(t, want, have)
	})
	t.Run("GitHub SAAS, no override", func(t *testing.T) {
		t.Parallel()
		url := giturl.Parse("username@github.com:git-town/docs.git")
		have := hosting.Detect(url, None[configdomain.HostingPlatform]())
		want := Some(configdomain.HostingPlatformGitHub)
		must.Eq(t, want, have)
	})
	t.Run("custom URL, override to GitHub", func(t *testing.T) {
		t.Parallel()
		url := giturl.Parse("username@custom.org:git-town/docs.git")
		have := hosting.Detect(url, Some(configdomain.HostingPlatformGitHub))
		want := Some(configdomain.HostingPlatformGitHub)
		must.Eq(t, want, have)
	})
	t.Run("GitLab SAAS, no override", func(t *testing.T) {
		t.Parallel()
		url := giturl.Parse("username@gitlab.com:git-town/docs.git")
		have := hosting.Detect(url, None[configdomain.HostingPlatform]())
		want := Some(configdomain.HostingPlatformGitLab)
		must.Eq(t, want, have)
	})
	t.Run("custom URL, override to GitLab", func(t *testing.T) {
		t.Parallel()
		url := giturl.Parse("username@custom.org:git-town/docs.git")
		have := hosting.Detect(url, Some(configdomain.HostingPlatformGitLab))
		want := Some(configdomain.HostingPlatformGitLab)
		must.Eq(t, want, have)
	})
}
