package mockproposals_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/test/mockproposals"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestMockProposals(t *testing.T) {
	t.Parallel()

	t.Run("FindBySourceAndTarget", func(t *testing.T) {
		t.Run("source and target match", func(t *testing.T) {
			t.Parallel()
			data1 := forgedomain.ProposalData{
				Number: 1,
				Source: "feature-branch",
				Target: "main",
			}
			data2 := forgedomain.ProposalData{
				Number: 2,
				Source: "other-branch",
				Target: "main",
			}
			proposals := mockproposals.MockProposals{data1, data2}
			have := proposals.FindBySourceAndTarget("feature-branch", "main")
			want := Some(data1)
			must.Eq(t, want, have)
		})

		t.Run("source matches but target does not", func(t *testing.T) {
			t.Parallel()
			proposals := mockproposals.MockProposals{
				{
					Number: 1,
					Source: "feature-branch",
					Target: "main",
					Title:  "Proposal 1",
					URL:    "https://example.com/pr/1",
				},
			}
			have := proposals.FindBySourceAndTarget("feature-branch", "develop")
			must.True(t, have.IsNone())
		})

		t.Run("target matches but source does not", func(t *testing.T) {
			t.Parallel()
			proposals := mockproposals.MockProposals{
				{
					Number: 1,
					Source: "feature-branch",
					Target: "main",
					Title:  "Proposal 1",
					URL:    "https://example.com/pr/1",
				},
			}
			have := proposals.FindBySourceAndTarget("other-branch", "main")
			must.True(t, have.IsNone())
		})

		t.Run("neither source nor target match", func(t *testing.T) {
			t.Parallel()
			proposals := mockproposals.MockProposals{
				{
					Number: 1,
					Source: "feature-branch",
					Target: "main",
					Title:  "Proposal 1",
					URL:    "https://example.com/pr/1",
				},
			}
			have := proposals.FindBySourceAndTarget("other-branch", "develop")
			must.True(t, have.IsNone())
		})

		t.Run("proposals slice is empty", func(t *testing.T) {
			t.Parallel()
			proposals := mockproposals.MockProposals{}
			have := proposals.FindBySourceAndTarget("feature-branch", "main")
			must.True(t, have.IsNone())
		})

		t.Run("multiple match", func(t *testing.T) {
			t.Parallel()
			data1 := forgedomain.ProposalData{
				Number: 1,
				Source: "feature-branch",
				Target: "main",
			}
			data2 := forgedomain.ProposalData{
				Number: 2,
				Source: "feature-branch",
				Target: "main",
			}
			proposals := mockproposals.MockProposals{data1, data2}
			have := proposals.FindBySourceAndTarget("feature-branch", "main")
			want := Some(data1)
			must.Eq(t, want, have)
		})
	})

	t.Run("FindById", func(t *testing.T) {
		t.Run("ID matches", func(t *testing.T) {
			t.Parallel()
			data1 := forgedomain.ProposalData{
				Number: 1,
				Source: "feature-branch",
				Target: "main",
			}
			data2 := forgedomain.ProposalData{
				Number: 2,
				Source: "other-branch",
				Target: "main",
			}
			proposals := mockproposals.MockProposals{data1, data2}
			have := proposals.FindById(2)
			want := MutableSome(&data2)
			must.Eq(t, want, have)
		})

		t.Run("ID does not match", func(t *testing.T) {
			t.Parallel()
			proposals := mockproposals.MockProposals{
				{
					Number: 1,
					Source: "feature-branch",
					Target: "main",
					Title:  "Proposal 1",
				},
			}
			have := proposals.FindById(999)
			must.True(t, have.IsNone())
		})

		t.Run("proposals slice is empty", func(t *testing.T) {
			t.Parallel()
			proposals := mockproposals.MockProposals{}
			have := proposals.FindById(1)
			must.True(t, have.IsNone())
		})
	})

	t.Run("Search", func(t *testing.T) {
		t.Run("returns all proposals matching source", func(t *testing.T) {
			t.Parallel()
			data1 := forgedomain.ProposalData{
				Number: 1,
				Source: "feature-branch",
				Target: "main",
			}
			data2 := forgedomain.ProposalData{
				Number: 2,
				Source: "feature-branch",
				Target: "develop",
			}
			data3 := forgedomain.ProposalData{
				Number: 3,
				Source: "other-branch",
				Target: "main",
			}
			proposals := mockproposals.MockProposals{data1, data2, data3}
			have := proposals.Search("feature-branch")
			want := []forgedomain.ProposalData{data1, data2}
			must.Eq(t, want, have)
		})

		t.Run("no proposals match source", func(t *testing.T) {
			t.Parallel()
			proposals := mockproposals.MockProposals{
				{
					Number: 1,
					Source: "feature-branch",
					Target: "main",
					Title:  "Proposal 1",
					URL:    "https://example.com/pr/1",
				},
			}
			have := proposals.Search("other-branch")
			must.Len(t, 0, have)
		})

		t.Run("proposals slice is empty", func(t *testing.T) {
			t.Parallel()
			proposals := mockproposals.MockProposals{}
			have := proposals.Search("feature-branch")
			must.Len(t, 0, have)
		})
	})
}
