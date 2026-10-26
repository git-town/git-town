package bitbucketcloud

import (
	"testing"

	"github.com/git-town/git-town/v24/internal/forge/forgedomain"
	"github.com/git-town/git-town/v24/internal/git/gitdomain"
	"github.com/shoenig/test/must"
)

func TestParsePullRequest(t *testing.T) {
	t.Parallel()

	newPullRequest := func() map[string]any {
		return map[string]any{
			"id":          float64(123),
			"title":       "title",
			"description": "body",
			"state":       "OPEN",
			"destination": map[string]any{
				"branch": map[string]any{
					"name": "main",
				},
			},
			"source": map[string]any{
				"branch": map[string]any{
					"name": "feature",
				},
			},
			"links": map[string]any{
				"html": map[string]any{
					"href": "https://bitbucket.org/org/repo/pull-requests/123",
				},
			},
			"close_source_branch": true,
			"draft":               false,
		}
	}

	t.Run("no reviewers", func(t *testing.T) {
		t.Parallel()
		give := newPullRequest()
		have, err := parsePullRequest(give)
		must.NoError(t, err)
		must.Eq(t, []string{}, have.Reviewers)
	})

	t.Run("reviewer without UUID", func(t *testing.T) {
		t.Parallel()
		give := newPullRequest()
		give["reviewers"] = []any{
			map[string]any{"account_id": "account-1"},
		}
		_, err := parsePullRequest(give)
		must.Error(t, err)
	})

	t.Run("reviewers", func(t *testing.T) {
		t.Parallel()
		give := newPullRequest()
		give["reviewers"] = []any{
			map[string]any{"uuid": "{reviewer-1}", "account_id": "account-1"},
			map[string]any{"uuid": "{reviewer-2}"},
		}
		have, err := parsePullRequest(give)
		must.NoError(t, err)
		want := forgedomain.BitbucketCloudProposalData{
			ProposalData: forgedomain.ProposalData{
				Active:       true,
				Body:         gitdomain.NewProposalBodyOpt("body"),
				MergeWithAPI: false,
				Number:       123,
				Source:       "feature",
				Target:       "main",
				Title:        "title",
				URL:          "https://bitbucket.org/org/repo/pull-requests/123",
			},
			CloseSourceBranch: true,
			Draft:             false,
			Reviewers:         []string{"{reviewer-1}", "{reviewer-2}"},
		}
		must.Eq(t, want, have)
	})
}
