package bitbucketcloud

import (
	"testing"

	"github.com/git-town/git-town/v24/internal/forge/forgedomain"
	"github.com/git-town/git-town/v24/internal/git/gitdomain"
	"github.com/shoenig/test/must"
)

func TestProposalBodyUpdateOptions(t *testing.T) {
	t.Parallel()

	connector := APIConnector{
		WebConnector: WebConnector{
			HostedRepoInfo: forgedomain.HostedRepoInfo{
				Organization: "org",
				Repository:   "repo",
			},
		},
	}
	proposalData := forgedomain.ProposalData{
		Number: 123,
		Source: "feature",
		Target: "main",
		Title:  "title",
		Body:   gitdomain.NewProposalBodyOpt("existing body"),
	}

	have := connector.proposalBodyUpdateOptions(proposalData, gitdomain.ProposalBody("updated body"))

	must.EqOp(t, "123", have.ID)
	must.EqOp(t, "org", have.Owner)
	must.EqOp(t, "repo", have.RepoSlug)
	must.EqOp(t, "title", have.Title)
	must.EqOp(t, "updated body", have.Description)
	must.EqOp(t, "", have.SourceBranch)
	must.EqOp(t, "", have.DestinationBranch)
	must.False(t, have.Draft)
	must.False(t, have.CloseSourceBranch)
	must.Len(t, 0, have.Reviewers)
	must.Len(t, 0, have.ReviewerAccountIDs)
}
