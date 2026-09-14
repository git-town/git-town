package forgedomain_test

import (
	"encoding/json"
	"testing"

	"github.com/git-town/git-town/v24/internal/forge/forgedomain"
	"github.com/git-town/git-town/v24/internal/git/gitdomain"
	"github.com/shoenig/test/must"
)

func TestBitbucketCloudProposalData(t *testing.T) {
	t.Parallel()

	t.Run("MarshalJSON", func(t *testing.T) {
		t.Parallel()
		data := forgedomain.BitbucketCloudProposalData{
			ProposalData: forgedomain.ProposalData{
				Active:       true,
				Body:         gitdomain.NewProposalBodyOpt("body"),
				MergeWithAPI: true,
				Number:       123,
				Source:       "source",
				Target:       "target",
				Title:        "title",
				URL:          "url",
			},
			CloseSourceBranch:  true,
			Draft:              true,
			Reviewers:          []string{"{reviewer-1}", "{reviewer-2}"},
			ReviewerAccountIDs: []string{"account-1", "account-2"},
		}
		serialized, err := json.MarshalIndent(data, "", "  ")
		must.NoError(t, err)
		want := `
{
  "Active": true,
  "Body": "body",
  "MergeWithAPI": true,
  "Number": 123,
  "Source": "source",
  "Target": "target",
  "Title": "title",
  "URL": "url",
  "CloseSourceBranch": true,
  "Draft": true,
  "Reviewers": [
    "{reviewer-1}",
    "{reviewer-2}"
  ],
  "ReviewerAccountIDs": [
    "account-1",
    "account-2"
  ]
}`[1:]
		must.EqOp(t, want, string(serialized))

		var data2 forgedomain.BitbucketCloudProposalData
		must.NoError(t, json.Unmarshal(serialized, &data2))
		must.Eq(t, data, data2)
		must.True(t, data2.CloseSourceBranch)
		must.True(t, data2.Draft)
		must.Eq(t, []string{"{reviewer-1}", "{reviewer-2}"}, data2.Reviewers)
		must.Eq(t, []string{"account-1", "account-2"}, data2.ReviewerAccountIDs)
	})
}
