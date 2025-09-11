package azuredevops_test

import (
	"testing"

	"github.com/git-town/git-town/v21/internal/forge/azuredevops"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/git/giturl"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestWebConnector(t *testing.T) {
	t.Parallel()
	t.Run("NewProposalURL", func(t *testing.T) {
		t.Parallel()
		url := giturl.Parse("git@ssh.dev.azure.com:v3/kevingoslar/tikibase/tikibase").GetOrPanic()
		connector := azuredevops.NewConnector(azuredevops.NewConnectorArgs{
			ProposalOverride: None[forgedomain.ProposalOverride](),
			RemoteURL:        url,
		})
		have := connector.NewProposalURL(forgedomain.CreateProposalArgs{
			Branch:         "feature",
			FrontendRunner: nil,
			MainBranch:     "main",
			ParentBranch:   "parent",
			ProposalBody:   Some(gitdomain.ProposalBody("body")),
			ProposalTitle:  Some(gitdomain.ProposalTitle("title")),
		})
		want := "https://dev.azure.com/kevingoslar/tikibase/_git/tikibase/pullrequestcreate?sourceRef=kg-test&targetRef=main"
		must.EqOp(t, want, have)
	})
}
