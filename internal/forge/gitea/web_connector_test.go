package gitea_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/forge/gitea"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/shoenig/test/must"
)

//nolint:paralleltest  // mocks HTTP
func TestGitea(t *testing.T) {
	t.Run("DefaultProposalMessage", func(t *testing.T) {
		t.Run("without body", func(t *testing.T) {
			give := forgedomain.ProposalData{
				Body:   None[gitdomain.ProposalBody](),
				Number: 123,
				Title:  "my title",
			}
			want := "my title (#123)"
			connector := gitea.WebConnector{}
			have := connector.DefaultProposalMessage(give)
			must.EqOp(t, want, have)
		})
		t.Run("with body", func(t *testing.T) {
			give := forgedomain.ProposalData{
				Body:   gitdomain.NewProposalBodyOpt("body"),
				Number: 123,
				Title:  "my title",
			}
			want := "my title (#123)\n\nbody"
			connector := gitea.WebConnector{}
			have := connector.DefaultProposalMessage(give)
			must.EqOp(t, want, have)
		})
	})

	t.Run("NewProposalURL", func(t *testing.T) {
		connector := gitea.WebConnector{
			HostedRepoInfo: forgedomain.HostedRepoInfo{
				Hostname:     "gitea.com",
				Organization: "org",
				Repository:   "repo",
			},
		}
		have := connector.NewProposalURL(forgedomain.CreateProposalArgs{
			Branch:        "feature",
			MainBranch:    "main",
			ParentBranch:  "parent",
			ProposalBody:  gitdomain.NewProposalBodyOpt("body"),
			ProposalTitle: Some(gitdomain.ProposalTitle("title")),
		})
		must.EqOp(t, "https://gitea.com/org/repo/compare/parent...feature", have)
	})

	t.Run("RepositoryURL", func(t *testing.T) {
		connector := gitea.WebConnector{
			HostedRepoInfo: forgedomain.HostedRepoInfo{
				Hostname:     "gitea.com",
				Organization: "org",
				Repository:   "repo",
			},
		}
		have := connector.RepositoryURL()
		must.EqOp(t, "https://gitea.com/org/repo", have)
	})
}
