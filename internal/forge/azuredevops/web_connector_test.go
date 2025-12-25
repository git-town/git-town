package azuredevops_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/forge/azuredevops"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/git/giturl"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestWebConnector(t *testing.T) {
	t.Parallel()
	url := giturl.Parse("git@ssh.dev.azure.com:v3/kevingoslar/tikibase/tikibase").GetOrPanic()
	connector := azuredevops.NewConnector(azuredevops.NewConnectorArgs{
		ProposalOverride: None[forgedomain.ProposalOverride](),
		RemoteURL:        url,
	})

	t.Run("NewProposalURL", func(t *testing.T) {
		t.Parallel()
		have := connector.NewProposalURL(forgedomain.CreateProposalArgs{
			Branch:         "feature",
			FrontendRunner: nil,
			MainBranch:     "main",
			ParentBranch:   "parent",
			ProposalBody:   gitdomain.NewProposalBodyOpt("body"),
			ProposalTitle:  Some(gitdomain.ProposalTitle("title")),
		})
		want := "https://dev.azure.com/kevingoslar/tikibase/_git/tikibase/pullrequestcreate?sourceRef=feature&targetRef=parent"
		must.EqOp(t, want, have)
	})

	t.Run("RepositoryURL", func(t *testing.T) {
		t.Parallel()
		have := connector.RepositoryURL()
		want := "https://dev.azure.com/kevingoslar/tikibase/_git/tikibase"
		must.EqOp(t, want, have)
	})
}
