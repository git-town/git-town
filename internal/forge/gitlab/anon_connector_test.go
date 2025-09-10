package gitlab_test

import (
	"testing"

	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/forge/gitlab"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestGitLabAnonConnector(t *testing.T) {
	t.Parallel()

	t.Run("DefaultProposalMessage", func(t *testing.T) {
		t.Parallel()
		t.Run("without body", func(t *testing.T) {
			t.Parallel()
			connector := gitlab.AnonConnector{}
			give := forgedomain.ProposalData{
				Number: 123,
				Title:  "my title",
			}
			have := connector.DefaultProposalMessage(give)
			want := "my title (!123)"
			must.EqOp(t, want, have)
		})
		t.Run("with body", func(t *testing.T) {
			t.Parallel()
			connector := gitlab.AnonConnector{}
			give := forgedomain.ProposalData{
				Number: 123,
				Title:  "my title",
				Body:   Some("body"),
			}
			have := connector.DefaultProposalMessage(give)
			want := "my title (!123)\n\nbody"
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
				want:   "https://gitlab.com/organization/repo/-/merge_requests/new?merge_request%5Bsource_branch%5D=feature&merge_request%5Btarget_branch%5D=main",
			},
			"stacked change": {
				branch: "feature-3",
				parent: "feature-2",
				want:   "https://gitlab.com/organization/repo/-/merge_requests/new?merge_request%5Bsource_branch%5D=feature-3&merge_request%5Btarget_branch%5D=feature-2",
			},
			"special characters in branch name": {
				branch: "feature-#",
				parent: "main",
				want:   "https://gitlab.com/organization/repo/-/merge_requests/new?merge_request%5Bsource_branch%5D=feature-%23&merge_request%5Btarget_branch%5D=main",
			},
			"proposal with title": {
				branch: "feature",
				parent: "main",
				title:  Some(gitdomain.ProposalTitle("my title")),
				body:   None[gitdomain.ProposalBody](),
				want:   "https://gitlab.com/organization/repo/-/merge_requests/new?merge_request%5Bsource_branch%5D=feature&merge_request%5Btarget_branch%5D=main&merge_request%5Btitle%5D=my+title",
			},
			"proposal with body": {
				branch: "feature",
				parent: "main",
				title:  None[gitdomain.ProposalTitle](),
				body:   Some(gitdomain.ProposalBody("my body")),
				want:   "https://gitlab.com/organization/repo/-/merge_requests/new?merge_request%5Bdescription%5D=my+body&merge_request%5Bsource_branch%5D=feature&merge_request%5Btarget_branch%5D=main",
			},
			"proposal with title and body": {
				branch: "feature",
				parent: "main",
				title:  Some(gitdomain.ProposalTitle("my title")),
				body:   Some(gitdomain.ProposalBody("my body")),
				want:   "https://gitlab.com/organization/repo/-/merge_requests/new?merge_request%5Bdescription%5D=my+body&merge_request%5Bsource_branch%5D=feature&merge_request%5Btarget_branch%5D=main&merge_request%5Btitle%5D=my+title",
			},
		}
		for name, tt := range tests {
			t.Run(name, func(t *testing.T) {
				t.Parallel()
				connector := gitlab.AnonConnector{}
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
}
