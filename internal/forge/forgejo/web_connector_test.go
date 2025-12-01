package forgejo_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/forge/forgejo"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestWebConnector(t *testing.T) {
	t.Parallel()

	t.Run("DefaultProposalMessage", func(t *testing.T) {
		t.Parallel()
		t.Run("with body", func(t *testing.T) {
			t.Parallel()
			give := forgedomain.ProposalData{
				Body:   gitdomain.NewProposalBodyOpt("body"),
				Number: 123,
				Title:  "my title",
			}
			want := "my title (#123)\n\nbody"
			connector := forgejo.WebConnector{}
			have := connector.DefaultProposalMessage(give)
			must.EqOp(t, want, have)
		})
		t.Run("without body", func(t *testing.T) {
			t.Parallel()
			give := forgedomain.ProposalData{
				Body:   None[gitdomain.ProposalBody](),
				Number: 123,
				Title:  "my title",
			}
			want := "my title (#123)"
			connector := forgejo.WebConnector{}
			have := connector.DefaultProposalMessage(give)
			must.EqOp(t, want, have)
		})
	})

	t.Run("NewProposalURL", func(t *testing.T) {
		t.Parallel()
		connector := forgejo.WebConnector{
			HostedRepoInfo: forgedomain.HostedRepoInfo{
				Hostname:     "codeberg.org",
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
		must.EqOp(t, "https://codeberg.org/org/repo/compare/parent...feature", have)
	})

	t.Run("RepositoryURL", func(t *testing.T) {
		t.Parallel()
		connector := forgejo.WebConnector{
			HostedRepoInfo: forgedomain.HostedRepoInfo{
				Hostname:     "codeberg.org",
				Organization: "org",
				Repository:   "repo",
			},
		}
		have := connector.RepositoryURL()
		must.EqOp(t, "https://codeberg.org/org/repo", have)
	})
}
