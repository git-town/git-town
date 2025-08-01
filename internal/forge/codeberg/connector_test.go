package codeberg_test

import (
	"testing"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2"
	"github.com/git-town/git-town/v21/internal/cli/print"
	"github.com/git-town/git-town/v21/internal/forge/codeberg"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/forge/github"
	"github.com/git-town/git-town/v21/internal/git/giturl"

	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestConnector(t *testing.T) {
	t.Parallel()

	t.Run("DefaultProposalMessage", func(t *testing.T) {
		t.Parallel()
		t.Run("with body", func(t *testing.T) {
			t.Parallel()
			give := forgedomain.ProposalData{
				Body:   Some("body"),
				Number: 123,
				Title:  "my title",
			}
			want := "my title (#123)\n\nbody"
			connector := codeberg.Connector{}
			have := connector.DefaultProposalMessage(give)
			must.EqOp(t, want, have)
		})
		t.Run("without body", func(t *testing.T) {
			t.Parallel()
			give := forgedomain.ProposalData{
				Body:   None[string](),
				Number: 123,
				Title:  "my title",
			}
			want := "my title (#123)"
			connector := codeberg.Connector{}
			have := connector.DefaultProposalMessage(give)
			must.EqOp(t, want, have)
		})
	})

	// THIS TEST CONNECTS TO AN EXTERNAL INTERNET HOST,
	// WHICH MAKES IT SLOW AND FLAKY.
	// DISABLE AS NEEDED TO DEBUG THE GITEA CONNECTOR.
	//
	// t.Run("NewProposalURL", func(t *testing.T) {
	// 	t.Parallel()
	// 	connector, err := codeberg.NewConnector(codeberg.NewConnectorArgs{
	// 		APIToken:  None[configdomain.CodebergToken](),
	// 		Log:       print.Logger{},
	// 		RemoteURL: giturl.Parse("git@codeberg.org:git-town/docs.git").GetOrPanic(),
	// 	})
	// 	must.NoError(t, err)
	// 	have, err := connector.NewProposalURL("feature", "parent", "", "", "")
	// 	must.NoError(t, err)
	// 	must.EqOp(t, "https://codeberg.org/git-town/docs/compare/parent...feature", have)
	// })

	// THIS TEST CONNECTS TO AN EXTERNAL INTERNET HOST,
	// WHICH MAKES IT SLOW AND FLAKY.
	// DISABLE AS NEEDED TO DEBUG THE GITEA CONNECTOR.
	//
	// t.Run("RepositoryURL", func(t *testing.T) {
	// 	t.Parallel()
	// 	connector, err := codeberg.NewConnector(codeberg.NewConnectorArgs{
	// 		APIToken:  None[configdomain.CodebergToken](),
	// 		Log:       print.Logger{},
	// 		RemoteURL: giturl.Parse("git@codeberg.org:git-town/docs.git").GetOrPanic(),
	// 	})
	// 	must.NoError(t, err)
	// 	have := connector.RepositoryURL()
	// 	must.EqOp(t, "https://codeberg.org/git-town/docs", have)
	// })
}

func TestFilterPullRequests(t *testing.T) {
	t.Parallel()
	give := []*forgejo.PullRequest{
		// matching branch
		{
			Head: &forgejo.PRBranchInfo{Name: "branch"},
			Base: &forgejo.PRBranchInfo{Name: "target"},
		},
		// branch with different name
		{
			Head: &forgejo.PRBranchInfo{Name: "other"},
			Base: &forgejo.PRBranchInfo{Name: "target"},
		},
		// branch with different target
		{
			Head: &forgejo.PRBranchInfo{Name: "branch"},
			Base: &forgejo.PRBranchInfo{Name: "other"},
		},
	}
	want := []*forgejo.PullRequest{
		{
			Head: &forgejo.PRBranchInfo{Name: "branch"},
			Base: &forgejo.PRBranchInfo{Name: "target"},
		},
	}
	have := codeberg.FilterPullRequests(give, "branch", "target")
	must.Eq(t, want, have)
}

func TestNewConnector(t *testing.T) {
	t.Parallel()

	// THIS TEST CONNECTS TO AN EXTERNAL INTERNET HOST,
	// WHICH MAKES IT SLOW AND FLAKY.
	// DISABLE AS NEEDED TO DEBUG THE GITEA CONNECTOR.
	//
	// t.Run("Codeberg SaaS", func(t *testing.T) {
	// 	t.Parallel()
	// 	have, err := codeberg.NewConnector(codeberg.NewConnectorArgs{
	// 		APIToken:  None[configdomain.CodebergToken](),
	// 		Log:       print.Logger{},
	// 		RemoteURL: giturl.Parse("git@codeberg.org:git-town/docs.git").GetOrPanic(),
	// 	})
	// 	must.NoError(t, err)
	// 	want := forgedomain.Data{
	// 		Hostname:     "codeberg.org",
	// 		Organization: "git-town",
	// 		Repository:   "docs",
	// 	}
	// 	must.EqOp(t, want, have.Data)
	// })

	t.Run("custom URL", func(t *testing.T) {
		t.Parallel()
		have, err := github.NewConnector(github.NewConnectorArgs{
			APIToken:  forgedomain.ParseGitHubToken("apiToken"),
			Log:       print.Logger{},
			RemoteURL: giturl.Parse("git@custom-url.com:git-town/docs.git").GetOrPanic(),
		})
		must.NoError(t, err)
		wantConfig := forgedomain.Data{
			Hostname:     "custom-url.com",
			Organization: "git-town",
			Repository:   "docs",
		}
		must.EqOp(t, wantConfig, have.Data)
	})
}
