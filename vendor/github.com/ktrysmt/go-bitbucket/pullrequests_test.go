package bitbucket

import (
	"encoding/json"
	"testing"

	"github.com/shoenig/test/must"
)

func TestBuildPullRequestBody(t *testing.T) {
	t.Parallel()

	t.Run("omits empty fields", func(t *testing.T) {
		t.Parallel()

		pullRequests := PullRequests{}

		have, err := pullRequests.buildPullRequestBody(&PullRequestsOptions{
			Owner:             "workspace",
			RepoSlug:          "repo",
			Title:             "title",
			Description:       "body",
			SourceBranch:      "feature",
			DestinationBranch: "main",
		})

		must.NoError(t, err)
		var payload map[string]any
		must.NoError(t, json.Unmarshal([]byte(have), &payload))
		must.Eq(t, "title", payload["title"])
		must.Eq(t, "body", payload["description"])
		_, hasReviewers := payload["reviewers"]
		must.False(t, hasReviewers)
		_, hasMessage := payload["message"]
		must.False(t, hasMessage)
		_, hasCloseSourceBranch := payload["close_source_branch"]
		must.False(t, hasCloseSourceBranch)
	})

	t.Run("includes reviewers when provided", func(t *testing.T) {
		t.Parallel()

		pullRequests := PullRequests{}

		have, err := pullRequests.buildPullRequestBody(&PullRequestsOptions{
			Reviewers: []string{"{reviewer-1}", "{reviewer-2}"},
		})

		must.NoError(t, err)
		var payload map[string]any
		must.NoError(t, json.Unmarshal([]byte(have), &payload))
		reviewers, hasReviewers := payload["reviewers"].([]any)
		must.True(t, hasReviewers)
		must.Len(t, 2, reviewers)
	})

	t.Run("prefers reviewer account ids when provided", func(t *testing.T) {
		t.Parallel()

		pullRequests := PullRequests{}

		have, err := pullRequests.buildPullRequestBody(&PullRequestsOptions{
			ReviewerAccountIDs: []string{"account-1", "account-2"},
		})

		must.NoError(t, err)
		var payload map[string]any
		must.NoError(t, json.Unmarshal([]byte(have), &payload))
		reviewers, hasReviewers := payload["reviewers"].([]any)
		must.True(t, hasReviewers)
		firstReviewer, ok := reviewers[0].(map[string]any)
		must.True(t, ok)
		must.Eq(t, "user", firstReviewer["type"])
		must.Eq(t, "account-1", firstReviewer["account_id"])
	})
}
