package confighelpers_test

import (
	"testing"

	"github.com/git-town/git-town/v15/internal/config/configdomain"
	"github.com/git-town/git-town/v15/internal/config/confighelpers"
	"github.com/git-town/git-town/v15/internal/git/giturl"
	. "github.com/git-town/git-town/v15/internal/gohacks/prelude"
	"github.com/shoenig/test/must"
)

func TestDetermineRemoteURL(t *testing.T) {
	t.Parallel()

	t.Run("SSH URL", func(t *testing.T) {
		t.Parallel()
		have, has := confighelpers.DetermineRemoteURL("git@github.com:git-town/docs.git", configdomain.ParseHostingOriginHostname("")).Get()
		must.True(t, has)
		want := giturl.Parts{
			Host: "github.com",
			Org:  "git-town",
			Repo: "docs",
			User: Some("git"),
		}
		must.Eq(t, want, have)
	})

	t.Run("HTTPS URL", func(t *testing.T) {
		t.Parallel()
		have, has := confighelpers.DetermineRemoteURL("https://github.com/git-town/docs.git", configdomain.ParseHostingOriginHostname("")).Get()
		must.True(t, has)
		want := giturl.Parts{
			Host: "github.com",
			Org:  "git-town",
			Repo: "docs",
			User: None[string](),
		}
		must.Eq(t, want, have)
	})

	t.Run("GitLab handbook repo on gitlab.com", func(t *testing.T) {
		t.Parallel()
		have, has := confighelpers.DetermineRemoteURL("git@gitlab.com:gitlab-com/www-gitlab-com.git", configdomain.ParseHostingOriginHostname("")).Get()
		must.True(t, has)
		want := giturl.Parts{
			Host: "gitlab.com",
			Org:  "gitlab-com",
			Repo: "www-gitlab-com",
			User: Some("git"),
		}
		must.Eq(t, want, have)
	})

	t.Run("GitLab repository inside a group", func(t *testing.T) {
		t.Parallel()
		have, has := confighelpers.DetermineRemoteURL("git@gitlab.com:gitlab-org/quality/triage-ops.git", configdomain.ParseHostingOriginHostname("")).Get()
		must.True(t, has)
		want := giturl.Parts{
			Host: "gitlab.com",
			Org:  "gitlab-org/quality",
			Repo: "triage-ops",
			User: Some("git"),
		}
		must.Eq(t, want, have)
	})

	t.Run("self-hosted GitLab server without URL override", func(t *testing.T) {
		t.Parallel()
		have, has := confighelpers.DetermineRemoteURL("git@self-hosted-gitlab.com:git-town/git-town.git", configdomain.ParseHostingOriginHostname("")).Get()
		must.True(t, has)
		want := giturl.Parts{
			Host: "self-hosted-gitlab.com",
			Org:  "git-town",
			Repo: "git-town",
			User: Some("git"),
		}
		must.Eq(t, want, have)
	})

	t.Run("self-hosted GitLab server with URL override", func(t *testing.T) {
		t.Parallel()
		have, has := confighelpers.DetermineRemoteURL("git@self-hosted-gitlab.com:git-town/git-town.git", configdomain.ParseHostingOriginHostname("override.com")).Get()
		must.True(t, has)
		want := giturl.Parts{
			Host: "override.com",
			Org:  "git-town",
			Repo: "git-town",
			User: Some("git"),
		}
		must.Eq(t, want, have)
	})

	t.Run("custom SSH identity with hostname override", func(t *testing.T) {
		t.Parallel()
		have, has := confighelpers.DetermineRemoteURL("git@my-ssh-identity.com:git-town/git-town.git", configdomain.ParseHostingOriginHostname("gitlab.com")).Get()
		must.True(t, has)
		want := giturl.Parts{
			Host: "gitlab.com",
			Org:  "git-town",
			Repo: "git-town",
			User: Some("git"),
		}
		must.Eq(t, want, have)
	})
}
