package forge_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/forge"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/giturl"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestDetect(t *testing.T) {
	t.Parallel()

	t.Run("Azure DevOps", func(t *testing.T) {
		t.Parallel()
		url, has := giturl.Parse("username@ssh.dev.azure.com:v3/kevingoslar/tikibase/tikibase").Get()
		must.True(t, has)
		have := forge.Detect(url, None[forgedomain.ForgeType]())
		want := Some(forgedomain.ForgeTypeAzureDevOps)
		must.Eq(t, want, have)
	})

	t.Run("BitBucket SAAS, no override", func(t *testing.T) {
		t.Parallel()
		url, has := giturl.Parse("username@bitbucket.org:git-town/docs.git").Get()
		must.True(t, has)
		have := forge.Detect(url, None[forgedomain.ForgeType]())
		want := Some(forgedomain.ForgeTypeBitbucket)
		must.Eq(t, want, have)
	})

	t.Run("GitHub SAAS, no override", func(t *testing.T) {
		t.Parallel()
		url, has := giturl.Parse("username@github.com:git-town/docs.git").Get()
		must.True(t, has)
		have := forge.Detect(url, None[forgedomain.ForgeType]())
		want := Some(forgedomain.ForgeTypeGitHub)
		must.Eq(t, want, have)
	})

	t.Run("GitLab SAAS, no override", func(t *testing.T) {
		t.Parallel()
		url, err := giturl.Parse("username@gitlab.com:git-town/docs.git").Get()
		must.True(t, err)
		have := forge.Detect(url, None[forgedomain.ForgeType]())
		want := Some(forgedomain.ForgeTypeGitLab)
		must.Eq(t, want, have)
	})

	t.Run("Gitea SAAS, no override", func(t *testing.T) {
		t.Parallel()
		url, has := giturl.Parse("username@gitea.com:git-town/docs.git").Get()
		must.True(t, has)
		have := forge.Detect(url, None[forgedomain.ForgeType]())
		want := Some(forgedomain.ForgeTypeGitea)
		must.Eq(t, want, have)
	})

	t.Run("custom URL, override to BitBucket", func(t *testing.T) {
		t.Parallel()
		url, has := giturl.Parse("username@custom.org:git-town/docs.git").Get()
		must.True(t, has)
		have := forge.Detect(url, Some(forgedomain.ForgeTypeBitbucket))
		want := Some(forgedomain.ForgeTypeBitbucket)
		must.Eq(t, want, have)
	})

	t.Run("custom URL, override to Bitbucket Data Center", func(t *testing.T) {
		t.Parallel()
		url, has := giturl.Parse("username@custom.org:git-town/docs.git").Get()
		must.True(t, has)
		have := forge.Detect(url, Some(forgedomain.ForgeTypeBitbucketDatacenter))
		want := Some(forgedomain.ForgeTypeBitbucketDatacenter)
		must.Eq(t, want, have)
	})

	t.Run("custom URL, override to GitHub", func(t *testing.T) {
		t.Parallel()
		url, has := giturl.Parse("username@custom.org:git-town/docs.git").Get()
		must.True(t, has)
		have := forge.Detect(url, Some(forgedomain.ForgeTypeGitHub))
		want := Some(forgedomain.ForgeTypeGitHub)
		must.Eq(t, want, have)
	})

	t.Run("custom URL, override to GitLab", func(t *testing.T) {
		t.Parallel()
		url, err := giturl.Parse("username@custom.org:git-town/docs.git").Get()
		must.True(t, err)
		have := forge.Detect(url, Some(forgedomain.ForgeTypeGitLab))
		want := Some(forgedomain.ForgeTypeGitLab)
		must.Eq(t, want, have)
	})

	t.Run("custom URL, override to Gitea", func(t *testing.T) {
		t.Parallel()
		url, has := giturl.Parse("username@custom.org:git-town/docs.git").Get()
		must.True(t, has)
		have := forge.Detect(url, Some(forgedomain.ForgeTypeGitea))
		want := Some(forgedomain.ForgeTypeGitea)
		must.Eq(t, want, have)
	})
}
