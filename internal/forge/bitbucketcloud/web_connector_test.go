package bitbucketcloud_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/forge/bitbucketcloud"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestWebConnector(t *testing.T) {
	t.Parallel()

	t.Run("NewProposalURL", func(t *testing.T) {
		t.Parallel()
		connector := bitbucketcloud.WebConnector{
			HostedRepoInfo: forgedomain.HostedRepoInfo{
				Hostname:     "bitbucket.org",
				Organization: "org",
				Repository:   "repo",
			},
		}
		have := connector.NewProposalURL(forgedomain.CreateProposalArgs{
			Branch:        "branch",
			MainBranch:    "main",
			ParentBranch:  "parent-branch",
			ProposalBody:  gitdomain.NewProposalBodyOpt("body"),
			ProposalTitle: Some(gitdomain.ProposalTitle("title")),
		})
		want := "https://bitbucket.org/org/repo/pull-requests/new?source=branch&dest=org%2Frepo%3Aparent-branch"
		must.EqOp(t, want, have)
	})
}
