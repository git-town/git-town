package forgejo_test

import (
	"testing"

	"github.com/git-town/git-town/v21/internal/cli/print"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/forge/github"
	"github.com/git-town/git-town/v21/internal/git/giturl"

	"github.com/shoenig/test/must"
)

func TestConnector(t *testing.T) {
	t.Parallel()

	// THIS TEST CONNECTS TO AN EXTERNAL INTERNET HOST,
	// WHICH MAKES IT SLOW AND FLAKY.
	// DISABLE AS NEEDED TO DEBUG THE GITEA CONNECTOR.
	//

	// THIS TEST CONNECTS TO AN EXTERNAL INTERNET HOST,
	// WHICH MAKES IT SLOW AND FLAKY.
	// DISABLE AS NEEDED TO DEBUG THE GITEA CONNECTOR.
	//
	// t.Run("RepositoryURL", func(t *testing.T) {
	// 	t.Parallel()
	// 	connector, err := forgejo.NewConnector(forgejo.NewConnectorArgs{
	// 		APIToken:  None[configdomain.ForgejoToken](),
	// 		Log:       print.Logger{},
	// 		RemoteURL: giturl.Parse("git@codeberg.org:git-town/docs.git").GetOrPanic(),
	// 	})
	// 	must.NoError(t, err)
	// 	have := connector.RepositoryURL()
	// 	must.EqOp(t, "https://codeberg.org/git-town/docs", have)
	// })
}

func TestNewConnector(t *testing.T) {
	t.Parallel()

	// THIS TEST CONNECTS TO AN EXTERNAL INTERNET HOST,
	// WHICH MAKES IT SLOW AND FLAKY.
	// DISABLE AS NEEDED TO DEBUG THE GITEA CONNECTOR.
	//
	// t.Run("Codeberg SaaS", func(t *testing.T) {
	// 	t.Parallel()
	// 	have, err := forgejo.NewConnector(forgejo.NewConnectorArgs{
	// 		APIToken:  None[configdomain.ForgejoToken](),
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
