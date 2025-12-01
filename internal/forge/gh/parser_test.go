package gh_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/forge/gh"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestParseJSONOutput(t *testing.T) {
	t.Parallel()

	t.Run("invalid JSON", func(t *testing.T) {
		t.Parallel()
		give := `[zonk`
		_, err := gh.ParseJSONOutput(give)
		must.Error(t, err)
	})

	t.Run("multiple results", func(t *testing.T) {
		t.Parallel()
		give := `
[
  {
    "baseRefName": "main",
    "body": "GitLab also provides a CLI app. This PR adds support for it similar to GitHub.\n",
    "headRefName": "kg-glab",
    "mergeable": "MERGEABLE",
    "number": 5079,
    "title": "glab connector type",
    "url": "https://github.com/git-town/git-town/pull/5079"
  },
  {
    "baseRefName": "main",
    "body": "Addresses:\r\n\r\nhttps://github.com/git-town/git-town/issues/3003",
    "headRefName": "support-pull-requests-comments",
    "mergeable": "UNKNOWN",
    "number": 4871,
    "title": "Feat: Display lineage / hierarchy for Proposals",
    "url": "https://github.com/git-town/git-town/pull/4871"
  }
]`
		_, err := gh.ParseJSONOutput(give)
		must.NoError(t, err)
	})

	t.Run("no results", func(t *testing.T) {
		t.Parallel()
		give := `[]`
		have, err := gh.ParseJSONOutput(give)
		must.NoError(t, err)
		must.Eq(t, []forgedomain.Proposal{}, have)
	})

	t.Run("single result", func(t *testing.T) {
		t.Parallel()
		give := `
[
  {
    "baseRefName": "main",
    "body": "GitLab also provides a CLI app. This PR adds support for it similar to GitHub.\n",
    "headRefName": "kg-glab",
    "mergeable": "MERGEABLE",
    "number": 5079,
    "title": "glab connector type",
    "url": "https://github.com/git-town/git-town/pull/5079"
  }
]`
		have, err := gh.ParseJSONOutput(give)
		must.NoError(t, err)
		want := []forgedomain.Proposal{
			{
				Data: forgedomain.ProposalData{
					Body:         Some(gitdomain.ProposalBody("GitLab also provides a CLI app. This PR adds support for it similar to GitHub.\n")),
					MergeWithAPI: true,
					Number:       5079,
					Source:       "kg-glab",
					Target:       "main",
					Title:        "glab connector type",
					URL:          "https://github.com/git-town/git-town/pull/5079",
				},
				ForgeType: forgedomain.ForgeTypeGitHub,
			},
		}
		must.Eq(t, want, have)
	})
}
