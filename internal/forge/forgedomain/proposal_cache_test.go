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

		t.Run("clears all cached results", func(t *testing.T) {
			t.Parallel()
			cache := &forgedomain.ProposalCache{}
			source := gitdomain.NewLocalBranchName("feature")
			target := gitdomain.NewLocalBranchName("main")
			proposal := forgedomain.Proposal{
				Data: forgedomain.ProposalData{
					Source: source,
					Target: target,
					Number: 123,
				},
				ForgeType: forgedomain.ForgeTypeGitHub,
			}
			cache.RegisterLookupResult(source, target, Some(proposal))
			cache.RegisterSearchResult(source, []forgedomain.Proposal{proposal})
			cache.Clear()
			got, knows := cache.Lookup(source, target)
			must.False(t, knows)
			must.True(t, got.IsNone())
			proposals, knows := cache.LookupSearch(source)
			must.False(t, knows)
			must.EqOp(t, 0, len(proposals))
		})
	})

	t.Run("RegisterLookupResult", func(t *testing.T) {
		t.Parallel()

		t.Run("registers new lookup result", func(t *testing.T) {
			t.Parallel()
			cache := &forgedomain.ProposalCache{}
			source := gitdomain.NewLocalBranchName("feature")
			target := gitdomain.NewLocalBranchName("main")
			proposal := forgedomain.Proposal{
				Data: forgedomain.ProposalData{
					Source: source,
					Target: target,
					Number: 123,
				},
				ForgeType: forgedomain.ForgeTypeGitHub,
			}
			cache.RegisterLookupResult(source, target, Some(proposal))
			got, knows := cache.Lookup(source, target)
			must.True(t, knows)
			must.True(t, got.IsSome())
		})

		t.Run("overwrites existing lookup result", func(t *testing.T) {
			t.Parallel()
			cache := &forgedomain.ProposalCache{}
			source := gitdomain.NewLocalBranchName("feature")
			target := gitdomain.NewLocalBranchName("main")
			proposal1 := forgedomain.Proposal{
				Data: forgedomain.ProposalData{
					Source: source,
					Target: target,
					Number: 123,
					Title:  "First PR",
				},
				ForgeType: forgedomain.ForgeTypeGitHub,
			}
			proposal2 := forgedomain.Proposal{
				Data: forgedomain.ProposalData{
					Source: source,
					Target: target,
					Number: 456,
					Title:  "Second PR",
				},
				ForgeType: forgedomain.ForgeTypeGitHub,
			}
			cache.RegisterLookupResult(source, target, Some(proposal1))
			cache.RegisterLookupResult(source, target, Some(proposal2))
			got, knows := cache.Lookup(source, target)
			must.True(t, knows)
			must.True(t, got.IsSome())
			gotProposal, _ := got.Get()
			must.EqOp(t, proposal2.Data.Data().Number, gotProposal.Data.Data().Number)
			must.EqOp(t, "Second PR", gotProposal.Data.Data().Title)
		})

		t.Run("can register multiple different source-target pairs", func(t *testing.T) {
			t.Parallel()
			cache := &forgedomain.ProposalCache{}
			source1 := gitdomain.NewLocalBranchName("feature1")
			source2 := gitdomain.NewLocalBranchName("feature2")
			target := gitdomain.NewLocalBranchName("main")
			proposal1 := forgedomain.Proposal{
				Data: forgedomain.ProposalData{
					Source: source1,
					Target: target,
					Number: 123,
				},
				ForgeType: forgedomain.ForgeTypeGitHub,
			}
			proposal2 := forgedomain.Proposal{
				Data: forgedomain.ProposalData{
					Source: source2,
					Target: target,
					Number: 456,
				},
				ForgeType: forgedomain.ForgeTypeGitHub,
			}
			cache.RegisterLookupResult(source1, target, Some(proposal1))
			cache.RegisterLookupResult(source2, target, Some(proposal2))
			got1, knows1 := cache.Lookup(source1, target)
			must.True(t, knows1)
			gotProposal1, _ := got1.Get()
			must.EqOp(t, proposal1.Data.Data().Number, gotProposal1.Data.Data().Number)
			got2, knows2 := cache.Lookup(source2, target)
			must.True(t, knows2)
			gotProposal2, _ := got2.Get()
			must.EqOp(t, proposal2.Data.Data().Number, gotProposal2.Data.Data().Number)
		})
	})

	t.Run("RegisterSearchResult", func(t *testing.T) {
		t.Parallel()

		t.Run("registers new search result", func(t *testing.T) {
			t.Parallel()
			cache := &forgedomain.ProposalCache{}
			source := gitdomain.NewLocalBranchName("feature")
			proposal := forgedomain.Proposal{
				Data: forgedomain.ProposalData{
					Source: source,
					Target: gitdomain.NewLocalBranchName("main"),
					Number: 123,
				},
				ForgeType: forgedomain.ForgeTypeGitHub,
			}
			cache.RegisterSearchResult(source, []forgedomain.Proposal{proposal})
			got, knows := cache.LookupSearch(source)
			must.True(t, knows)
			must.EqOp(t, 1, len(got))
		})

		t.Run("overwrites existing search result", func(t *testing.T) {
			t.Parallel()
			cache := &forgedomain.ProposalCache{}
			source := gitdomain.NewLocalBranchName("feature")
			proposal1 := forgedomain.Proposal{
				Data: forgedomain.ProposalData{
					Source: source,
					Target: gitdomain.NewLocalBranchName("main"),
					Number: 123,
				},
				ForgeType: forgedomain.ForgeTypeGitHub,
			}
			proposal2 := forgedomain.Proposal{
				Data: forgedomain.ProposalData{
					Source: source,
					Target: gitdomain.NewLocalBranchName("develop"),
					Number: 456,
				},
				ForgeType: forgedomain.ForgeTypeGitHub,
			}
			cache.RegisterSearchResult(source, []forgedomain.Proposal{proposal1})
			cache.RegisterSearchResult(source, []forgedomain.Proposal{proposal2})
			got, knows := cache.LookupSearch(source)
			must.True(t, knows)
			must.EqOp(t, 1, len(got))
			must.EqOp(t, proposal2.Data.Data().Number, got[0].Data.Data().Number)
		})

		t.Run("can register multiple different sources", func(t *testing.T) {
			t.Parallel()
			cache := &forgedomain.ProposalCache{}
			source1 := gitdomain.NewLocalBranchName("feature1")
			source2 := gitdomain.NewLocalBranchName("feature2")
			proposal1 := forgedomain.Proposal{
				Data: forgedomain.ProposalData{
					Source: source1,
					Target: gitdomain.NewLocalBranchName("main"),
					Number: 123,
				},
				ForgeType: forgedomain.ForgeTypeGitHub,
			}
			proposal2 := forgedomain.Proposal{
				Data: forgedomain.ProposalData{
					Source: source2,
					Target: gitdomain.NewLocalBranchName("main"),
					Number: 456,
				},
				ForgeType: forgedomain.ForgeTypeGitHub,
			}
			cache.RegisterSearchResult(source1, []forgedomain.Proposal{proposal1})
			cache.RegisterSearchResult(source2, []forgedomain.Proposal{proposal2})
			got1, knows1 := cache.LookupSearch(source1)
			must.True(t, knows1)
			must.EqOp(t, 1, len(got1))
			must.EqOp(t, proposal1.Data.Data().Number, got1[0].Data.Data().Number)
			got2, knows2 := cache.LookupSearch(source2)
			must.True(t, knows2)
			must.EqOp(t, 1, len(got2))
			must.EqOp(t, proposal2.Data.Data().Number, got2[0].Data.Data().Number)
		})

		t.Run("lookup result takes precedence over search result when both exist", func(t *testing.T) {
			t.Parallel()
			cache := &forgedomain.ProposalCache{}
			source := gitdomain.NewLocalBranchName("feature")
			target := gitdomain.NewLocalBranchName("main")
			lookupProposal := forgedomain.Proposal{
				Data: forgedomain.ProposalData{
					Source: source,
					Target: target,
					Number: 123,
				},
				ForgeType: forgedomain.ForgeTypeGitHub,
			}
			searchProposal := forgedomain.Proposal{
				Data: forgedomain.ProposalData{
					Source: source,
					Target: target,
					Number: 456,
				},
				ForgeType: forgedomain.ForgeTypeGitHub,
			}
			cache.RegisterLookupResult(source, target, Some(lookupProposal))
			cache.RegisterSearchResult(source, []forgedomain.Proposal{searchProposal})
			// Lookup result takes precedence because it's checked first in the switch statement
			got, knows := cache.Lookup(source, target)
			must.True(t, knows)
			must.True(t, got.IsSome())
			gotProposal, _ := got.Get()
			must.EqOp(t, lookupProposal.Data.Data().Number, gotProposal.Data.Data().Number)
		})
	})
}
