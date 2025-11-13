package forgedomain_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestProposalCache(t *testing.T) {
	t.Parallel()

	t.Run("Lookup", func(t *testing.T) {
		t.Parallel()

		t.Run("empty cache returns unknown", func(t *testing.T) {
			t.Parallel()
			cache := &forgedomain.ProposalCache{}
			have, has := cache.Lookup("source", "target")
			must.False(t, has)
			must.True(t, have.IsNone())
		})

		t.Run("cache contains existing lookup result", func(t *testing.T) {
			t.Parallel()

			t.Run("matching lookup result", func(t *testing.T) {
				t.Parallel()

				t.Run("has proposal", func(t *testing.T) {
					t.Parallel()
					cache := &forgedomain.ProposalCache{}
					giveProposal := forgedomain.Proposal{
						Data: forgedomain.ProposalData{
							Source: "source",
							Target: "target",
							Number: 123,
							Title:  "Test PR",
						},
						ForgeType: forgedomain.ForgeTypeGitHub,
					}
					cache.RegisterLookupResult("source", "target", Some(giveProposal))
					lookupResult, has := cache.Lookup("source", "target")
					must.True(t, has)
					haveProposal, has := lookupResult.Get()
					must.True(t, has)
					must.EqOp(t, 123, haveProposal.Data.Data().Number)
					must.EqOp(t, "Test PR", haveProposal.Data.Data().Title)
				})

				t.Run("without proposal", func(t *testing.T) {
					t.Parallel()
					cache := &forgedomain.ProposalCache{}
					cache.RegisterLookupResult(source, target, None[forgedomain.Proposal]())
					lookupResult, has := cache.Lookup(source, target)
					must.True(t, has)
					must.True(t, lookupResult.IsNone())
				})
			})

			t.Run("mismatching lookup result", func(t *testing.T) {
				t.Parallel()

				t.Run("different source", func(t *testing.T) {
					t.Parallel()
					cache := &forgedomain.ProposalCache{}
					giveProposal := forgedomain.Proposal{
						Data: forgedomain.ProposalData{
							Source: "other",
							Target: "target",
							Number: 123,
							Title:  "Test PR",
						},
						ForgeType: forgedomain.ForgeTypeGitHub,
					}
					cache.RegisterLookupResult("other", "target", Some(giveProposal))
					lookupResult, has := cache.Lookup("source", "target")
					must.False(t, has)
				})

				t.Run("different target", func(t *testing.T) {
					t.Parallel()
					cache := &forgedomain.ProposalCache{}
					giveProposal := forgedomain.Proposal{
						Data: forgedomain.ProposalData{
							Source: "source",
							Target: "other",
							Number: 123,
							Title:  "Test PR",
						},
						ForgeType: forgedomain.ForgeTypeGitHub,
					}
					cache.RegisterLookupResult("source", "other", Some(giveProposal))
					lookupResult, has := cache.Lookup("source", "target")
					must.False(t, has)
				})
			})
		})

		t.Run("cache contains existing search result", func(t *testing.T) {
			t.Parallel()

			t.Run("matching search result that contains the target", func(t *testing.T) {
				t.Parallel()
				cache := &forgedomain.ProposalCache{}
				proposal1 := forgedomain.Proposal{
					Data: forgedomain.ProposalData{
						Source: "source",
						Target: "target",
						Number: 123,
						Title:  "Test PR",
					},
					ForgeType: forgedomain.ForgeTypeGitHub,
				}
				cache.RegisterSearchResult("source", []forgedomain.Proposal{proposal1})
				lookupResult, has := cache.Lookup("source", "target")
				must.True(t, has)
				haveProposal, has := lookupResult.Get()
				must.True(t, has)
				must.EqOp(t, 123, haveProposal.Data.Data().Number)
			})

			t.Run("matching search result that does not contain the target", func(t *testing.T) {
				t.Parallel()
				cache := &forgedomain.ProposalCache{}
				proposal1 := forgedomain.Proposal{
					Data: forgedomain.ProposalData{
						Source: "source",
						Target: "other",
						Number: 123,
						Title:  "Test PR",
					},
					ForgeType: forgedomain.ForgeTypeGitHub,
				}
				cache.RegisterSearchResult("source", []forgedomain.Proposal{proposal1})
				lookupResult, has := cache.Lookup("source", "target")
				must.True(t, has)
				must.True(t, lookupResult.IsNone())
			})

			t.Run("matching search result that contains no proposals", func(t *testing.T) {
				t.Parallel()
				cache := &forgedomain.ProposalCache{}
				cache.RegisterSearchResult("source", []forgedomain.Proposal{})
				lookupResult, has := cache.Lookup("source", "target")
				must.True(t, has)
				must.True(t, lookupResult.IsNone())
			})
		})

		t.Run("cache contains both lookup and search result", func(t *testing.T) {
			t.Parallel()
			cache := &forgedomain.ProposalCache{}
			lookupProposal := forgedomain.Proposal{
				Data: forgedomain.ProposalData{
					Source: "source",
					Target: "target",
					Number: 123,
					Title:  "Lookup PR",
				},
				ForgeType: forgedomain.ForgeTypeGitHub,
			}
			searchProposal := forgedomain.Proposal{
				Data: forgedomain.ProposalData{
					Source: "source",
					Target: "target",
					Number: 456,
					Title:  "Search PR",
				},
				ForgeType: forgedomain.ForgeTypeGitHub,
			}
			// Register lookup first so it's checked before search result in the loop
			cache.RegisterLookupResult("source", "target", Some(lookupProposal))
			cache.RegisterSearchResult("source", []forgedomain.Proposal{searchProposal})
			result, has := cache.Lookup("source", "target")
			must.True(t, has)
			haveProposal, has := result.Get()
			must.True(t, has)
			must.EqOp(t, "Lookup PR", haveProposal.Data.Data().Title)
		})
	})

	t.Run("LookupSearch", func(t *testing.T) {
		t.Parallel()

		t.Run("empty cache returns unknown", func(t *testing.T) {
			t.Parallel()
			cache := &forgedomain.ProposalCache{}
			source := gitdomain.NewLocalBranchName("feature")
			_, knows := cache.LookupSearch(source)
			must.False(t, knows)
		})

		t.Run("contains matching search result", func(t *testing.T) {
			t.Parallel()
			cache := &forgedomain.ProposalCache{}
			proposal1 := forgedomain.Proposal{
				Data: forgedomain.ProposalData{
					Source: "source",
					Target: "target",
					Number: 123,
				},
				ForgeType: forgedomain.ForgeTypeGitHub,
			}
			proposal2 := forgedomain.Proposal{
				Data: forgedomain.ProposalData{
					Source: "source",
					Target: "target",
					Number: 456,
				},
				ForgeType: forgedomain.ForgeTypeGitHub,
			}
			giveProposals := []forgedomain.Proposal{proposal1, proposal2}
			cache.RegisterSearchResult("source", giveProposals)
			haveProposals, has := cache.LookupSearch("source")
			must.True(t, has)
			must.EqOp(t, 2, len(haveProposals))
			must.EqOp(t, 123, haveProposals[0].Data.Data().Number)
			must.EqOp(t, 456, haveProposals[1].Data.Data().Number)
		})

		t.Run("contains mismatching search result", func(t *testing.T) {
			t.Parallel()
			cache := &forgedomain.ProposalCache{}
			cache.RegisterSearchResult("other", []forgedomain.Proposal{})
			_, knows := cache.LookupSearch("source")
			must.False(t, knows)
		})

		t.Run("ignores lookup results", func(t *testing.T) {
			t.Parallel()
			cache := &forgedomain.ProposalCache{}
			proposal := forgedomain.Proposal{
				Data: forgedomain.ProposalData{
					Source: "source",
					Target: "target",
					Number: 123,
				},
				ForgeType: forgedomain.ForgeTypeGitHub,
			}
			cache.RegisterLookupResult("source", "target", Some(proposal))
			_, knows := cache.LookupSearch("source")
			must.False(t, knows)
		})
	})

	t.Run("Clear", func(t *testing.T) {
		t.Parallel()

		t.Run("removes all cached results", func(t *testing.T) {
			t.Parallel()
			cache := &forgedomain.ProposalCache{}
			proposal := forgedomain.Proposal{
				Data: forgedomain.ProposalData{
					Source: "source",
					Target: "target",
					Number: 123,
				},
				ForgeType: forgedomain.ForgeTypeGitHub,
			}
			cache.RegisterLookupResult("source", "target", Some(proposal))
			cache.RegisterSearchResult("source", []forgedomain.Proposal{proposal})
			cache.Clear()
			_, has := cache.Lookup("source", "target")
			must.False(t, has)
			_, has = cache.LookupSearch("source")
			must.False(t, has)
		})
	})

	t.Run("RegisterLookupResult", func(t *testing.T) {
		t.Parallel()

		t.Run("registers new lookup result", func(t *testing.T) {
			t.Parallel()
			cache := &forgedomain.ProposalCache{}
			proposal := forgedomain.Proposal{
				Data: forgedomain.ProposalData{
					Source: "source",
					Target: "target",
					Number: 123,
				},
				ForgeType: forgedomain.ForgeTypeGitHub,
			}
			cache.RegisterLookupResult("source", "target", Some(proposal))
			got, knows := cache.Lookup("source", "target")
			must.True(t, knows)
			must.True(t, got.IsSome())
		})

		t.Run("overwrites an existing lookup result", func(t *testing.T) {
			t.Parallel()
			cache := &forgedomain.ProposalCache{}
			proposal1 := forgedomain.Proposal{
				Data: forgedomain.ProposalData{
					Source: "source",
					Target: "target",
					Number: 123,
					Title:  "First PR",
				},
				ForgeType: forgedomain.ForgeTypeGitHub,
			}
			proposal2 := forgedomain.Proposal{
				Data: forgedomain.ProposalData{
					Source: "source",
					Target: "target",
					Number: 456,
					Title:  "Second PR",
				},
				ForgeType: forgedomain.ForgeTypeGitHub,
			}
			cache.RegisterLookupResult("source", "target", Some(proposal1))
			cache.RegisterLookupResult("source", "target", Some(proposal2))
			result, has := cache.Lookup("source", "target")
			must.True(t, has)
			haveProposal, has := result.Get()
			must.True(t, has)
			must.EqOp(t, "Second PR", haveProposal.Data.Data().Title)
		})

		t.Run("registers multiple different source-target pairs", func(t *testing.T) {
			t.Parallel()
			cache := &forgedomain.ProposalCache{}
			proposal1 := forgedomain.Proposal{
				Data: forgedomain.ProposalData{
					Source: "source1",
					Target: "target",
					Number: 123,
				},
				ForgeType: forgedomain.ForgeTypeGitHub,
			}
			proposal2 := forgedomain.Proposal{
				Data: forgedomain.ProposalData{
					Source: "source2",
					Target: "target",
					Number: 456,
				},
				ForgeType: forgedomain.ForgeTypeGitHub,
			}
			cache.RegisterLookupResult("source1", "target", Some(proposal1))
			cache.RegisterLookupResult("source2", "target", Some(proposal2))
			result1, has1 := cache.Lookup("source1", "target")
			must.True(t, has1)
			haveProposal1, has := result1.Get()
			must.True(t, has)
			must.EqOp(t, 123, haveProposal1.Data.Data().Number)
			result2, has := cache.Lookup("source2", "target")
			must.True(t, has)
			haveProposal2, has := result2.Get()
			must.True(t, has)
			must.EqOp(t, 456, haveProposal2.Data.Data().Number)
		})
	})

	t.Run("RegisterSearchResult", func(t *testing.T) {
		t.Parallel()

		t.Run("registers new search result", func(t *testing.T) {
			t.Parallel()
			cache := &forgedomain.ProposalCache{}
			proposal := forgedomain.Proposal{
				Data: forgedomain.ProposalData{
					Source: "source",
					Target: "target",
					Number: 123,
				},
				ForgeType: forgedomain.ForgeTypeGitHub,
			}
			cache.RegisterSearchResult("source", []forgedomain.Proposal{proposal})
			result, has := cache.LookupSearch("source")
			must.True(t, has)
			must.EqOp(t, 1, len(result))
			must.EqOp(t, 123, result[0].Data.Data().Number)
		})

		t.Run("overwrites existing search result", func(t *testing.T) {
			t.Parallel()
			cache := &forgedomain.ProposalCache{}
			proposal1 := forgedomain.Proposal{
				Data: forgedomain.ProposalData{
					Source: "source",
					Target: "target",
					Number: 123,
				},
				ForgeType: forgedomain.ForgeTypeGitHub,
			}
			proposal2 := forgedomain.Proposal{
				Data: forgedomain.ProposalData{
					Source: "source",
					Target: "target",
					Number: 456,
				},
				ForgeType: forgedomain.ForgeTypeGitHub,
			}
			cache.RegisterSearchResult("source", []forgedomain.Proposal{proposal1})
			cache.RegisterSearchResult("source", []forgedomain.Proposal{proposal2})
			haveProposals, has := cache.LookupSearch("source")
			must.True(t, has)
			must.EqOp(t, 1, len(haveProposals))
			must.EqOp(t, 456, haveProposals[0].Data.Data().Number)
		})

		t.Run("registers multiple different sources", func(t *testing.T) {
			t.Parallel()
			cache := &forgedomain.ProposalCache{}
			proposal1 := forgedomain.Proposal{
				Data: forgedomain.ProposalData{
					Source: "source1",
					Target: "target",
					Number: 123,
				},
				ForgeType: forgedomain.ForgeTypeGitHub,
			}
			proposal2 := forgedomain.Proposal{
				Data: forgedomain.ProposalData{
					Source: "source2",
					Target: "target",
					Number: 456,
				},
				ForgeType: forgedomain.ForgeTypeGitHub,
			}
			cache.RegisterSearchResult("source1", []forgedomain.Proposal{proposal1})
			cache.RegisterSearchResult("source2", []forgedomain.Proposal{proposal2})
			haveProposals1, has := cache.LookupSearch("source1")
			must.True(t, has)
			must.EqOp(t, 1, len(haveProposals1))
			must.EqOp(t, 123, haveProposals1[0].Data.Data().Number)
			haveProposals2, has := cache.LookupSearch("source2")
			must.True(t, has)
			must.EqOp(t, 1, len(haveProposals2))
			must.EqOp(t, 456, haveProposals2[0].Data.Data().Number)
		})
	})
}
