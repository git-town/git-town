package forgedomain_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestAPICache(t *testing.T) {
	t.Parallel()

	t.Run("Clear", func(t *testing.T) {
		t.Parallel()

		t.Run("removes the proposal with the given number from the cache", func(t *testing.T) {
			t.Parallel()
			cache := &forgedomain.APICache{}
			proposal1 := forgedomain.Proposal{
				Data: forgedomain.ProposalData{
					Number: 1,
					Source: "source",
					Target: "target",
				},
			}
			proposal2 := forgedomain.Proposal{
				Data: forgedomain.ProposalData{
					Number: 2,
					Source: "other",
					Target: "target",
				},
			}
			cache.RegisterLookupResult("source", "target", Some(proposal1))
			cache.RegisterLookupResult("other", "target", Some(proposal2))
			cache.Clear(proposal1.Data.Data().Number)
			haveOpt, certain := cache.Lookup("source", "target")
			must.False(t, certain)
			_, has := haveOpt.Get()
			must.False(t, has)
			haveOpt, certain = cache.Lookup("other", "target")
			must.True(t, certain)
			have, has := haveOpt.Get()
			must.True(t, has)
			must.Eq(t, proposal2, have)
		})
	})

	t.Run("Lookup", func(t *testing.T) {
		t.Parallel()

		t.Run("empty cache returns unknown", func(t *testing.T) {
			t.Parallel()
			cache := &forgedomain.APICache{}
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
					cache := &forgedomain.APICache{}
					giveProposal := forgedomain.Proposal{
						Data: forgedomain.ProposalData{
							Source: "source",
							Target: "target",
							Title:  "Test PR",
						},
					}
					cache.RegisterLookupResult("source", "target", Some(giveProposal))
					lookupResult, has := cache.Lookup("source", "target")
					must.True(t, has)
					haveProposal, has := lookupResult.Get()
					must.True(t, has)
					must.EqOp(t, "Test PR", haveProposal.Data.Data().Title)
				})

				t.Run("without proposal", func(t *testing.T) {
					t.Parallel()
					cache := &forgedomain.APICache{}
					cache.RegisterLookupResult("source", "target", None[forgedomain.Proposal]())
					lookupResult, has := cache.Lookup("source", "target")
					must.True(t, has)
					must.True(t, lookupResult.IsNone())
				})
			})

			t.Run("mismatching lookup result", func(t *testing.T) {
				t.Parallel()

				t.Run("different source", func(t *testing.T) {
					t.Parallel()
					cache := &forgedomain.APICache{}
					giveProposal := forgedomain.Proposal{
						Data: forgedomain.ProposalData{
							Source: "other",
							Target: "target",
						},
					}
					cache.RegisterLookupResult("other", "target", Some(giveProposal))
					_, has := cache.Lookup("source", "target")
					must.False(t, has)
				})

				t.Run("different target", func(t *testing.T) {
					t.Parallel()
					cache := &forgedomain.APICache{}
					giveProposal := forgedomain.Proposal{
						Data: forgedomain.ProposalData{
							Source: "source",
							Target: "other",
						},
					}
					cache.RegisterLookupResult("source", "other", Some(giveProposal))
					_, has := cache.Lookup("source", "target")
					must.False(t, has)
				})
			})
		})

		t.Run("cache contains existing search result", func(t *testing.T) {
			t.Parallel()

			t.Run("matching search result that contains the target", func(t *testing.T) {
				t.Parallel()
				cache := &forgedomain.APICache{}
				proposal1 := forgedomain.Proposal{
					Data: forgedomain.ProposalData{
						Source: "source",
						Target: "target",
						Title:  "Test PR",
					},
				}
				cache.RegisterSearchResult("source", []forgedomain.Proposal{proposal1})
				lookupResult, has := cache.Lookup("source", "target")
				must.True(t, has)
				haveProposal, has := lookupResult.Get()
				must.True(t, has)
				must.EqOp(t, "Test PR", haveProposal.Data.Data().Title)
			})

			t.Run("matching search result that does not contain the target", func(t *testing.T) {
				t.Parallel()
				cache := &forgedomain.APICache{}
				proposal1 := forgedomain.Proposal{
					Data: forgedomain.ProposalData{
						Source: "source",
						Target: "other",
					},
				}
				cache.RegisterSearchResult("source", []forgedomain.Proposal{proposal1})
				lookupResult, has := cache.Lookup("source", "target")
				must.True(t, has)
				must.True(t, lookupResult.IsNone())
			})

			t.Run("matching search result that contains no proposals", func(t *testing.T) {
				t.Parallel()
				cache := &forgedomain.APICache{}
				cache.RegisterSearchResult("source", []forgedomain.Proposal{})
				lookupResult, has := cache.Lookup("source", "target")
				must.True(t, has)
				must.True(t, lookupResult.IsNone())
			})
		})

		t.Run("cache contains both lookup and search result", func(t *testing.T) {
			t.Parallel()
			cache := &forgedomain.APICache{}
			lookupProposal := forgedomain.Proposal{
				Data: forgedomain.ProposalData{
					Source: "source",
					Target: "target",
					Title:  "Lookup PR",
				},
			}
			searchProposal := forgedomain.Proposal{
				Data: forgedomain.ProposalData{
					Source: "source",
					Target: "target",
					Title:  "Search PR",
				},
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
			cache := &forgedomain.APICache{}
			source := gitdomain.NewLocalBranchName("feature")
			_, knows := cache.LookupSearch(source).Get()
			must.False(t, knows)
		})

		t.Run("contains matching search result", func(t *testing.T) {
			t.Parallel()
			cache := &forgedomain.APICache{}
			proposal1 := forgedomain.Proposal{
				Data: forgedomain.ProposalData{
					Source: "source",
					Target: "target",
					Title:  "PR 1",
				},
			}
			proposal2 := forgedomain.Proposal{
				Data: forgedomain.ProposalData{
					Source: "source",
					Target: "target",
					Title:  "PR 2",
				},
			}
			giveProposals := []forgedomain.Proposal{proposal1, proposal2}
			cache.RegisterSearchResult("source", giveProposals)
			haveProposals, has := cache.LookupSearch("source").Get()
			must.True(t, has)
			must.EqOp(t, 2, len(haveProposals))
			must.EqOp(t, "PR 1", haveProposals[0].Data.Data().Title)
			must.EqOp(t, "PR 2", haveProposals[1].Data.Data().Title)
		})

		t.Run("contains mismatching search result", func(t *testing.T) {
			t.Parallel()
			cache := &forgedomain.APICache{}
			cache.RegisterSearchResult("other", []forgedomain.Proposal{})
			_, knows := cache.LookupSearch("source").Get()
			must.False(t, knows)
		})

		t.Run("ignores lookup results", func(t *testing.T) {
			t.Parallel()
			cache := &forgedomain.APICache{}
			proposal := forgedomain.Proposal{
				Data: forgedomain.ProposalData{
					Source: "source",
					Target: "target",
				},
			}
			cache.RegisterLookupResult("source", "target", Some(proposal))
			_, knows := cache.LookupSearch("source").Get()
			must.False(t, knows)
		})
	})

	t.Run("RegisterLookupResult", func(t *testing.T) {
		t.Parallel()

		t.Run("registers new lookup result", func(t *testing.T) {
			t.Parallel()
			cache := &forgedomain.APICache{}
			proposal := forgedomain.Proposal{
				Data: forgedomain.ProposalData{
					Source: "source",
					Target: "target",
				},
			}
			cache.RegisterLookupResult("source", "target", Some(proposal))
			got, knows := cache.Lookup("source", "target")
			must.True(t, knows)
			must.True(t, got.IsSome())
		})

		t.Run("overwrites an existing lookup result", func(t *testing.T) {
			t.Parallel()
			cache := &forgedomain.APICache{}
			proposal1 := forgedomain.Proposal{
				Data: forgedomain.ProposalData{
					Source: "source",
					Target: "target",
					Title:  "First PR",
				},
			}
			proposal2 := forgedomain.Proposal{
				Data: forgedomain.ProposalData{
					Source: "source",
					Target: "target",
					Title:  "Second PR",
				},
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
			cache := &forgedomain.APICache{}
			proposal1 := forgedomain.Proposal{
				Data: forgedomain.ProposalData{
					Source: "source1",
					Target: "target",
					Title:  "First PR",
				},
			}
			proposal2 := forgedomain.Proposal{
				Data: forgedomain.ProposalData{
					Source: "source2",
					Target: "target",
					Title:  "Second PR",
				},
			}
			cache.RegisterLookupResult("source1", "target", Some(proposal1))
			cache.RegisterLookupResult("source2", "target", Some(proposal2))
			result1, has1 := cache.Lookup("source1", "target")
			must.True(t, has1)
			haveProposal1, has := result1.Get()
			must.True(t, has)
			must.EqOp(t, "First PR", haveProposal1.Data.Data().Title)
			result2, has := cache.Lookup("source2", "target")
			must.True(t, has)
			haveProposal2, has := result2.Get()
			must.True(t, has)
			must.EqOp(t, "Second PR", haveProposal2.Data.Data().Title)
		})
	})

	t.Run("RegisterSearchResult", func(t *testing.T) {
		t.Parallel()

		t.Run("registers new search result", func(t *testing.T) {
			t.Parallel()
			cache := &forgedomain.APICache{}
			proposal := forgedomain.Proposal{
				Data: forgedomain.ProposalData{
					Source: "source",
					Target: "target",
					Title:  "PR 1",
				},
			}
			cache.RegisterSearchResult("source", []forgedomain.Proposal{proposal})
			result, has := cache.LookupSearch("source").Get()
			must.True(t, has)
			must.EqOp(t, 1, len(result))
			must.EqOp(t, "PR 1", result[0].Data.Data().Title)
		})

		t.Run("overwrites existing search result", func(t *testing.T) {
			t.Parallel()
			cache := &forgedomain.APICache{}
			proposal1 := forgedomain.Proposal{
				Data: forgedomain.ProposalData{
					Source: "source",
					Target: "target",
					Title:  "First PR",
				},
			}
			proposal2 := forgedomain.Proposal{
				Data: forgedomain.ProposalData{
					Source: "source",
					Target: "target",
					Title:  "Second PR",
				},
			}
			cache.RegisterSearchResult("source", []forgedomain.Proposal{proposal1})
			cache.RegisterSearchResult("source", []forgedomain.Proposal{proposal2})
			haveProposals, has := cache.LookupSearch("source").Get()
			must.True(t, has)
			must.EqOp(t, 1, len(haveProposals))
			must.EqOp(t, "Second PR", haveProposals[0].Data.Data().Title)
		})

		t.Run("registers multiple different sources", func(t *testing.T) {
			t.Parallel()
			cache := &forgedomain.APICache{}
			proposal1 := forgedomain.Proposal{
				Data: forgedomain.ProposalData{
					Source: "source1",
					Target: "target",
					Title:  "First PR",
				},
			}
			proposal2 := forgedomain.Proposal{
				Data: forgedomain.ProposalData{
					Source: "source2",
					Target: "target",
					Title:  "Second PR",
				},
			}
			cache.RegisterSearchResult("source1", []forgedomain.Proposal{proposal1})
			cache.RegisterSearchResult("source2", []forgedomain.Proposal{proposal2})
			haveProposals1, has := cache.LookupSearch("source1").Get()
			must.True(t, has)
			must.EqOp(t, 1, len(haveProposals1))
			must.EqOp(t, "First PR", haveProposals1[0].Data.Data().Title)
			haveProposals2, has := cache.LookupSearch("source2").Get()
			must.True(t, has)
			must.EqOp(t, 1, len(haveProposals2))
			must.EqOp(t, "Second PR", haveProposals2[0].Data.Data().Title)
		})
	})
}
