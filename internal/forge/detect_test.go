package forge_test

import (
	"testing"

	"github.com/git-town/git-town/v19/internal/config/configdomain"
	"github.com/git-town/git-town/v19/internal/forge"
	"github.com/git-town/git-town/v19/internal/git/giturl"
	. "github.com/git-town/git-town/v19/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestDetect(t *testing.T) {
	t.Parallel()

	t.Run("BitBucket SAAS, no override", func(t *testing.T) {
		t.Parallel()
		url, has := giturl.Parse("username@bitbucket.org:git-town/docs.git").Get()
		must.True(t, has)
		have := forge.Detect(url, None[configdomain.ForgeType]())
		want := Some(configdomain.ForgeTypeBitbucket)
		must.Eq(t, want, have)
	})

	t.Run("GitHub SAAS, no override", func(t *testing.T) {
		t.Parallel()
		url, has := giturl.Parse("username@github.com:git-town/docs.git").Get()
		must.True(t, has)
		have := forge.Detect(url, None[configdomain.ForgeType]())
		want := Some(configdomain.ForgeTypeGitHub)
		must.Eq(t, want, have)
	})

	t.Run("GitLab SAAS, no override", func(t *testing.T) {
		t.Parallel()
		url, err := giturl.Parse("username@gitlab.com:git-town/docs.git").Get()
		must.True(t, err)
		have := forge.Detect(url, None[configdomain.ForgeType]())
		want := Some(configdomain.ForgeTypeGitLab)
		must.Eq(t, want, have)
	})

	t.Run("Gitea SAAS, no override", func(t *testing.T) {
		t.Parallel()
		url, has := giturl.Parse("username@gitea.com:git-town/docs.git").Get()
		must.True(t, has)
		have := forge.Detect(url, None[configdomain.ForgeType]())
		want := Some(configdomain.ForgeTypeGitea)
		must.Eq(t, want, have)
	})

	t.Run("custom URL, override to BitBucket", func(t *testing.T) {
		t.Parallel()
		url, has := giturl.Parse("username@custom.org:git-town/docs.git").Get()
		must.True(t, has)
		have := forge.Detect(url, Some(configdomain.ForgeTypeBitbucket))
		want := Some(configdomain.ForgeTypeBitbucket)
		must.Eq(t, want, have)
	})

	t.Run("custom URL, override to Bitbucket Data Center", func(t *testing.T) {
		t.Parallel()
		url, has := giturl.Parse("username@custom.org:git-town/docs.git").Get()
		must.True(t, has)
		have := forge.Detect(url, Some(configdomain.ForgeTypeBitbucketDatacenter))
		want := Some(configdomain.ForgeTypeBitbucketDatacenter)
		must.Eq(t, want, have)
	})

	t.Run("custom URL, override to GitHub", func(t *testing.T) {
		t.Parallel()
		url, has := giturl.Parse("username@custom.org:git-town/docs.git").Get()
		must.True(t, has)
		have := forge.Detect(url, Some(configdomain.ForgeTypeGitHub))
		want := Some(configdomain.ForgeTypeGitHub)
		must.Eq(t, want, have)
	})

	t.Run("custom URL, override to GitLab", func(t *testing.T) {
		t.Parallel()
		url, err := giturl.Parse("username@custom.org:git-town/docs.git").Get()
		must.True(t, err)
		have := forge.Detect(url, Some(configdomain.ForgeTypeGitLab))
		want := Some(configdomain.ForgeTypeGitLab)
		must.Eq(t, want, have)
	})

	t.Run("custom URL, override to Gitea", func(t *testing.T) {
		t.Parallel()
		url, has := giturl.Parse("username@custom.org:git-town/docs.git").Get()
		must.True(t, has)
		have := forge.Detect(url, Some(configdomain.ForgeTypeGitea))
		want := Some(configdomain.ForgeTypeGitea)
		must.Eq(t, want, have)
	})
}
