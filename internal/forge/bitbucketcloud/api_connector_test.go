package bitbucketcloud

import (
	"testing"

	"github.com/git-town/git-town/v24/internal/forge/forgedomain"
	"github.com/git-town/git-town/v24/internal/git/gitdomain"
	"github.com/shoenig/test/must"
)

func TestProposalUpdateOptions(t *testing.T) {
	t.Parallel()

	connector := APIConnector{
		WebConnector: WebConnector{
			HostedRepoInfo: forgedomain.HostedRepoInfo{
				Organization: "org",
				Repository:   "repo",
			},
		},
	}

	t.Run("no body", func(t *testing.T) {
		t.Parallel()
		data := forgedomain.BitbucketCloudProposalData{
			ProposalData: forgedomain.ProposalData{
				Number: 123,
				Source: "feature",
				Target: "main",
				Title:  "title",
				Body:   gitdomain.NewProposalBodyOpt(""),
			},
			CloseSourceBranch: false,
			Draft:             false,
			Reviewers:         []string{},
		}
		have := connector.proposalUpdateOptions(data)
		must.EqOp(t, "", have.Description)
		must.Len(t, 0, have.Reviewers)
	})

	t.Run("sends the complete state of the proposal", func(t *testing.T) {
		t.Parallel()
		data := forgedomain.BitbucketCloudProposalData{
			ProposalData: forgedomain.ProposalData{
				Number: 123,
				Source: "feature",
				Target: "main",
				Title:  "title",
				Body:   gitdomain.NewProposalBodyOpt("body"),
			},
			CloseSourceBranch: true,
			Draft:             true,
			Reviewers:         []string{"{reviewer-1}", "{reviewer-2}"},
		}
		have := connector.proposalUpdateOptions(data)
		must.EqOp(t, "123", have.ID)
		must.EqOp(t, "org", have.Owner)
		must.EqOp(t, "repo", have.RepoSlug)
		must.EqOp(t, "feature", have.SourceBranch)
		must.EqOp(t, "main", have.DestinationBranch)
		must.EqOp(t, "title", have.Title)
		must.EqOp(t, "body", have.Description)
		must.True(t, have.Draft)
		must.True(t, have.CloseSourceBranch)
		must.Eq(t, []string{"{reviewer-1}", "{reviewer-2}"}, have.Reviewers)
	})
}
