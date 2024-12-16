package hosting_test

import (
	"testing"

	"github.com/git-town/git-town/v17/internal/config/configdomain"
	"github.com/git-town/git-town/v17/internal/git/giturl"
	"github.com/git-town/git-town/v17/internal/hosting"
	. "github.com/git-town/git-town/v17/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestDetect(t *testing.T) {
	t.Parallel()

	t.Run("BitBucket SAAS, no override", func(t *testing.T) {
		t.Parallel()
		url, has := giturl.Parse("username@bitbucket.org:git-town/docs.git").Get()
		must.True(t, has)
		have := hosting.Detect(url, None[configdomain.HostingPlatform]())
		want := Some(configdomain.HostingPlatformBitbucket)
		must.Eq(t, want, have)
	})

	t.Run("custom URL, override to BitBucket", func(t *testing.T) {
		t.Parallel()
		url, has := giturl.Parse("username@custom.org:git-town/docs.git").Get()
		must.True(t, has)
		have := hosting.Detect(url, Some(configdomain.HostingPlatformBitbucket))
		want := Some(configdomain.HostingPlatformBitbucket)
		must.Eq(t, want, have)
	})

	t.Run("custom URL, override to BitBucket Datacenter", func(t *testing.T) {
		t.Parallel()
		url, has := giturl.Parse("username@custom.org:git-town/docs.git").Get()
		must.True(t, has)
		have := hosting.Detect(url, Some(configdomain.HostingPlatformBitbucketDatacenter))
		want := Some(configdomain.HostingPlatformBitbucketDatacenter)
		must.Eq(t, want, have)
	})

	t.Run("Gitea SAAS, no override", func(t *testing.T) {
		t.Parallel()
		url, has := giturl.Parse("username@gitea.com:git-town/docs.git").Get()
		must.True(t, has)
		have := hosting.Detect(url, None[configdomain.HostingPlatform]())
		want := Some(configdomain.HostingPlatformGitea)
		must.Eq(t, want, have)
	})

	t.Run("custom URL, override to Gitea", func(t *testing.T) {
		t.Parallel()
		url, has := giturl.Parse("username@custom.org:git-town/docs.git").Get()
		must.True(t, has)
		have := hosting.Detect(url, Some(configdomain.HostingPlatformGitea))
		want := Some(configdomain.HostingPlatformGitea)
		must.Eq(t, want, have)
	})

	t.Run("GitHub SAAS, no override", func(t *testing.T) {
		t.Parallel()
		url, has := giturl.Parse("username@github.com:git-town/docs.git").Get()
		must.True(t, has)
		have := hosting.Detect(url, None[configdomain.HostingPlatform]())
		want := Some(configdomain.HostingPlatformGitHub)
		must.Eq(t, want, have)
	})

	t.Run("custom URL, override to GitHub", func(t *testing.T) {
		t.Parallel()
		url, has := giturl.Parse("username@custom.org:git-town/docs.git").Get()
		must.True(t, has)
		have := hosting.Detect(url, Some(configdomain.HostingPlatformGitHub))
		want := Some(configdomain.HostingPlatformGitHub)
		must.Eq(t, want, have)
	})

	t.Run("GitLab SAAS, no override", func(t *testing.T) {
		t.Parallel()
		url, err := giturl.Parse("username@gitlab.com:git-town/docs.git").Get()
		must.True(t, err)
		have := hosting.Detect(url, None[configdomain.HostingPlatform]())
		want := Some(configdomain.HostingPlatformGitLab)
		must.Eq(t, want, have)
	})

	t.Run("custom URL, override to GitLab", func(t *testing.T) {
		t.Parallel()
		url, err := giturl.Parse("username@custom.org:git-town/docs.git").Get()
		must.True(t, err)
		have := hosting.Detect(url, Some(configdomain.HostingPlatformGitLab))
		want := Some(configdomain.HostingPlatformGitLab)
		must.Eq(t, want, have)
	})
}
