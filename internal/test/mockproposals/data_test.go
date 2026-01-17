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
	t.Run("DeleteByID", func(t *testing.T) {
		t.Parallel()

		t.Run("delete middle one", func(t *testing.T) {
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
			data3 := forgedomain.ProposalData{
				Number: 3,
				Source: "third-branch",
				Target: "main",
			}
			proposals := mockproposals.MockProposals{data1, data2, data3}
			proposals.DeleteByID(2)
			want := mockproposals.MockProposals{data1, data3}
			must.Eq(t, want, proposals)
		})

		t.Run("no match", func(t *testing.T) {
			t.Parallel()
			data1 := forgedomain.ProposalData{
				Number: 1,
				Source: "feature-branch",
				Target: "main",
				Title:  "Proposal 1",
			}
			proposals := mockproposals.MockProposals{data1}
			proposals.DeleteByID(999)
			want := mockproposals.MockProposals{data1}
			must.Eq(t, want, proposals)
		})

		t.Run("no proposals", func(t *testing.T) {
			t.Parallel()
			proposals := mockproposals.MockProposals{}
			proposals.DeleteByID(1)
			want := mockproposals.MockProposals{}
			must.Eq(t, want, proposals)
		})
	})

	t.Run("FindByID", func(t *testing.T) {
		t.Parallel()
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
			have := proposals.FindByID(2)
			want := Some(data2)
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
			have := proposals.FindByID(999)
			must.True(t, have.IsNone())
		})

		t.Run("no proposals", func(t *testing.T) {
			t.Parallel()
			proposals := mockproposals.MockProposals{}
			have := proposals.FindByID(1)
			must.True(t, have.IsNone())
		})
	})

	t.Run("FindBySource", func(t *testing.T) {
		t.Parallel()
		t.Run("returns all proposals matching source", func(t *testing.T) {
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
			data3 := forgedomain.ProposalData{
				Number: 3,
				Source: "feature-branch",
				Target: "develop",
			}
			proposals := mockproposals.MockProposals{data1, data2, data3}
			have := proposals.FindBySource("feature-branch")
			want := []forgedomain.ProposalData{data1, data3}
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
			have := proposals.FindBySource("other-branch")
			must.Len(t, 0, have)
		})

		t.Run("proposals slice is empty", func(t *testing.T) {
			t.Parallel()
			proposals := mockproposals.MockProposals{}
			have := proposals.FindBySource("feature-branch")
			must.Len(t, 0, have)
		})
	})

	t.Run("FindBySourceAndTarget", func(t *testing.T) {
		t.Parallel()
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
}
