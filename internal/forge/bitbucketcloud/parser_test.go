package bitbucketcloud

import (
	"testing"

	"github.com/shoenig/test/must"
)

func TestParsePullRequest(t *testing.T) {
	t.Parallel()

	t.Run("preserves reviewer UUIDs", func(t *testing.T) {
		t.Parallel()

		give := map[string]any{
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
			"reviewers": []any{
				map[string]any{"uuid": "{reviewer-1}", "account_id": "account-1"},
				map[string]any{"uuid": "{reviewer-2}", "account_id": "account-2"},
			},
		}

		have, err := parsePullRequest(give)

		must.NoError(t, err)
		must.Eq(t, []string{"{reviewer-1}", "{reviewer-2}"}, have.Reviewers)
		must.Eq(t, []string{"account-1", "account-2"}, have.ReviewerAccountIDs)
	})

	t.Run("keeps all reviewer uuids when some account ids are missing", func(t *testing.T) {
		t.Parallel()

		give := map[string]any{
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
			"reviewers": []any{
				map[string]any{"uuid": "{reviewer-1}", "account_id": "account-1"},
				map[string]any{"uuid": "{reviewer-2}"},
			},
		}

		have, err := parsePullRequest(give)

		must.NoError(t, err)
		must.Eq(t, []string{"{reviewer-1}", "{reviewer-2}"}, have.Reviewers)
		must.Eq(t, []string{"account-1"}, have.ReviewerAccountIDs)
	})
}
