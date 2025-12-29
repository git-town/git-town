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
					URL:    "https://example.com/pr/1",
				},
			}
			have := proposals.FindById(999)
			must.True(t, have.IsNone())
		})

		t.Run("returns MutableNone when proposals slice is empty", func(t *testing.T) {
			t.Parallel()
			proposals := mockproposals.MockProposals{}
			have := proposals.FindById(1)
			must.True(t, have.IsNone())
		})

		t.Run("finds proposal with ID 0", func(t *testing.T) {
			t.Parallel()
			proposals := mockproposals.MockProposals{
				{
					Number: 0,
					Source: "feature-branch",
					Target: "main",
					Title:  "Proposal 0",
					URL:    "https://example.com/pr/0",
				},
				{
					Number: 1,
					Source: "other-branch",
					Target: "main",
					Title:  "Proposal 1",
					URL:    "https://example.com/pr/1",
				},
			}
			have := proposals.FindById(0)
			must.True(t, have.IsSome())
			value := have.GetOrPanic()
			must.EqOp(t, 0, value.Number)
		})
	})
}

func TestMockProposals_Search(t *testing.T) {
	t.Parallel()

	t.Run("returns all proposals matching source", func(t *testing.T) {
		t.Parallel()
		proposals := mockproposals.MockProposals{
			{
				Number: 1,
				Source: "feature-branch",
				Target: "main",
				Title:  "Proposal 1",
				URL:    "https://example.com/pr/1",
			},
			{
				Number: 2,
				Source: "feature-branch",
				Target: "develop",
				Title:  "Proposal 2",
				URL:    "https://example.com/pr/2",
			},
			{
				Number: 3,
				Source: "other-branch",
				Target: "main",
				Title:  "Proposal 3",
				URL:    "https://example.com/pr/3",
			},
		}
		have := proposals.Search("feature-branch")
		must.Len(t, 2, have)
		must.EqOp(t, 1, have[0].Number)
		must.EqOp(t, 2, have[1].Number)
	})

	t.Run("returns empty slice when no proposals match source", func(t *testing.T) {
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

	t.Run("returns empty slice when proposals slice is empty", func(t *testing.T) {
		t.Parallel()
		proposals := mockproposals.MockProposals{}
		have := proposals.Search("feature-branch")
		must.Len(t, 0, have)
	})

	t.Run("returns single proposal when only one matches", func(t *testing.T) {
		t.Parallel()
		proposals := mockproposals.MockProposals{
			{
				Number: 1,
				Source: "feature-branch",
				Target: "main",
				Title:  "Proposal 1",
				URL:    "https://example.com/pr/1",
			},
			{
				Number: 2,
				Source: "other-branch",
				Target: "main",
				Title:  "Proposal 2",
				URL:    "https://example.com/pr/2",
			},
		}
		have := proposals.Search("feature-branch")
		must.Len(t, 1, have)
		must.EqOp(t, 1, have[0].Number)
	})

	t.Run("returns all proposals when all match source", func(t *testing.T) {
		t.Parallel()
		proposals := mockproposals.MockProposals{
			{
				Number: 1,
				Source: "feature-branch",
				Target: "main",
				Title:  "Proposal 1",
				URL:    "https://example.com/pr/1",
			},
			{
				Number: 2,
				Source: "feature-branch",
				Target: "develop",
				Title:  "Proposal 2",
				URL:    "https://example.com/pr/2",
			},
		}
		have := proposals.Search("feature-branch")
		must.Len(t, 2, have)
		must.EqOp(t, 1, have[0].Number)
		must.EqOp(t, 2, have[1].Number)
	})
}
