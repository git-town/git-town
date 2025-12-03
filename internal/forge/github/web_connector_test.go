package github_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/forge/github"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestWebConnector(t *testing.T) {
	t.Parallel()

	t.Run("DefaultProposalMessage", func(t *testing.T) {
		t.Parallel()
		t.Run("without body", func(t *testing.T) {
			t.Parallel()
			connector := github.WebConnector{}
			give := forgedomain.ProposalData{
				Number: 123,
				Title:  "my title",
			}
			have := connector.DefaultProposalMessage(give)
			want := "my title (#123)"
			must.EqOp(t, want, have)
		})
		t.Run("with body", func(t *testing.T) {
			t.Parallel()
			connector := github.WebConnector{}
			give := forgedomain.ProposalData{
				Number: 123,
				Title:  "my title",
				Body:   gitdomain.NewProposalBodyOpt("body"),
			}
			have := connector.DefaultProposalMessage(give)
			want := "my title (#123)\n\nbody"
			must.EqOp(t, want, have)
		})
	})

	t.Run("NewProposalURL", func(t *testing.T) {
		t.Parallel()
		tests := map[string]struct {
			branch gitdomain.LocalBranchName
			parent gitdomain.LocalBranchName
			title  Option[gitdomain.ProposalTitle]
			body   Option[gitdomain.ProposalBody]
			want   string
		}{
			"top-level branch": {
				branch: "feature",
				parent: "main",
				title:  None[gitdomain.ProposalTitle](),
				body:   None[gitdomain.ProposalBody](),
				want:   "https://github.com/organization/repo/compare/feature?expand=1",
			},
			"stacked change": {
				branch: "feature-3",
				parent: "feature-2",
				title:  None[gitdomain.ProposalTitle](),
				body:   None[gitdomain.ProposalBody](),
				want:   "https://github.com/organization/repo/compare/feature-2...feature-3?expand=1",
			},
			"special characters in branch name": {
				branch: "feature-#",
				parent: "main",
				title:  None[gitdomain.ProposalTitle](),
				body:   None[gitdomain.ProposalBody](),
				want:   "https://github.com/organization/repo/compare/feature-%23?expand=1",
			},
			"provide title and body": {
				branch: "feature-#",
				parent: "main",
				title:  Some(gitdomain.ProposalTitle("my title")),
				body:   gitdomain.NewProposalBodyOpt("my body"),
				want:   "https://github.com/organization/repo/compare/feature-%23?expand=1&title=my+title&body=my+body",
			},
			"provide title only": {
				branch: "feature-#",
				parent: "main",
				title:  Some(gitdomain.ProposalTitle("my title")),
				body:   None[gitdomain.ProposalBody](),
				want:   "https://github.com/organization/repo/compare/feature-%23?expand=1&title=my+title",
			},
			"provide body only": {
				branch: "feature-#",
				parent: "main",
				title:  None[gitdomain.ProposalTitle](),
				body:   gitdomain.NewProposalBodyOpt("my body"),
				want:   "https://github.com/organization/repo/compare/feature-%23?expand=1&body=my+body",
			},
		}
		for name, tt := range tests {
			t.Run(name, func(t *testing.T) {
				t.Parallel()
				connector := github.WebConnector{
					HostedRepoInfo: forgedomain.HostedRepoInfo{
						Hostname:     "github.com",
						Organization: "organization",
						Repository:   "repo",
					},
				}
				have := connector.NewProposalURL(forgedomain.CreateProposalArgs{
					Branch:        tt.branch,
					MainBranch:    "main",
					ParentBranch:  tt.parent,
					ProposalBody:  tt.body,
					ProposalTitle: tt.title,
				})
				must.EqOp(t, tt.want, have)
			})
		}
	})

	t.Run("RepositoryURL", func(t *testing.T) {
		t.Parallel()
		connector := github.WebConnector{
			HostedRepoInfo: forgedomain.HostedRepoInfo{
				Hostname:     "github.com",
				Organization: "organization",
				Repository:   "repo",
			},
		}
		have := connector.RepositoryURL()
		want := "https://github.com/organization/repo"
		must.EqOp(t, want, have)
	})
}
