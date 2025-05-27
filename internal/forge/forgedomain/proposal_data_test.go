package forgedomain_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestBitbucketCloudProposalData(t *testing.T) {
	t.Parallel()

	t.Run("MarshalJSON", func(t *testing.T) {
		t.Parallel()
		data := forgedomain.BitbucketCloudProposalData{
			ProposalData: forgedomain.ProposalData{
				Body:         Some("body"),
				MergeWithAPI: true,
				Number:       123,
				Source:       "source",
				Target:       "target",
				Title:        "title",
				URL:          "url",
			},
			CloseSourceBranch: true,
			Draft:             true,
		}
		serialized, err := json.MarshalIndent(data, "", "  ")
		must.NoError(t, err)
		fmt.Println(string(serialized))
		want := `
{
  "Body": "body",
  "MergeWithAPI": true,
  "Number": 123,
  "Source": "source",
  "Target": "target",
  "Title": "title",
  "URL": "url",
  "CloseSourceBranch": true,
  "Draft": true
}`[1:]
		must.EqOp(t, want, string(serialized))

		var data2 forgedomain.BitbucketCloudProposalData
		must.NoError(t, json.Unmarshal(serialized, &data2))
		must.Eq(t, data, data2)
		must.True(t, data2.CloseSourceBranch)
		must.True(t, data2.Draft)
	})
}
