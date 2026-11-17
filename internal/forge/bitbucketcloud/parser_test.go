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
			"close_source_branch": true,
			"description":         "body",
			"destination": map[string]any{
				"branch": map[string]any{
					"name": "main",
				},
			},
			"draft": false,
			"id":    float64(123),
			"links": map[string]any{
				"html": map[string]any{
					"href": "https://bitbucket.org/org/repo/pull-requests/123",
				},
			},
			"source": map[string]any{
				"branch": map[string]any{
					"name": "feature",
				},
			},
			"state": "OPEN",
			"title": "title",
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
			map[string]any{"account_id": "account-1", "uuid": "{reviewer-1}"},
			map[string]any{"uuid": "{reviewer-2}"},
		}
		have, err := parsePullRequest(give)
		must.NoError(t, err)
		want := forgedomain.BitbucketCloudProposalData{
			Active:            true,
			Body:              gitdomain.NewProposalBodyOpt("body"),
			CloseSourceBranch: true,
			Draft:             false,
			MergeWithAPI:      false,
			Number:            123,
			Reviewers:         []string{"{reviewer-1}", "{reviewer-2}"},
			Source:            "feature",
			Target:            "main",
			Title:             "title",
			URL:               "https://bitbucket.org/org/repo/pull-requests/123",
		}
		must.Eq(t, want, have)
	})
}
