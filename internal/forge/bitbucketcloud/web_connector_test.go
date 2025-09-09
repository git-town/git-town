package bitbucketcloud_test

import (
	"testing"

	"github.com/git-town/git-town/v21/internal/forge/bitbucketcloud"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestWebConnector(t *testing.T) {
	t.Run("NewProposalURL", func(t *testing.T) {
		t.Parallel()
		connector := bitbucketcloud.WebConnector{
			Data: forgedomain.Data{
				Hostname:     "bitbucket.org",
				Organization: "org",
				Repository:   "repo",
			},
		}
		have := connector.NewProposalURL(forgedomain.CreateProposalArgs{
			Branch:        "branch",
			MainBranch:    "main",
			ParentBranch:  "parent-branch",
			ProposalBody:  None[gitdomain.ProposalBody](),
			ProposalTitle: None[gitdomain.ProposalTitle](),
		})
		want := "https://bitbucket.org/org/repo/pull-requests/new?source=branch&dest=org%2Frepo%3Aparent-branch"
		must.EqOp(t, want, have)
	})
}
