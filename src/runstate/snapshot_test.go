package runstate_test

import (
	"fmt"
	"testing"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/runstate"
	"github.com/stretchr/testify/assert"
)

func TestSnapshot(t *testing.T) {
	t.Parallel()
	t.Run("Diff", func(t *testing.T) {
		t.Parallel()
		t.Run("branches added", func(t *testing.T) {
			t.Parallel()
			before := runstate.Snapshot{
				Branches: domain.BranchInfos{
					domain.BranchInfo{
						Name:       domain.NewLocalBranchName("branch-2"),
						InitialSHA: domain.NewSHA("222222"),
					},
				},
			}
			after := runstate.Snapshot{
				Branches: domain.BranchInfos{
					domain.BranchInfo{
						Name:       domain.NewLocalBranchName("branch-1"),
						InitialSHA: domain.NewSHA("111111"),
					},
					domain.BranchInfo{
						Name:       domain.NewLocalBranchName("branch-2"),
						InitialSHA: domain.NewSHA("222222"),
					},
					domain.BranchInfo{
						Name:       domain.NewLocalBranchName("branch-3"),
						InitialSHA: domain.NewSHA("333333"),
					},
				},
			}
			have := after.Diff(before)
			want := runstate.Diff{
				BranchesUpdated: map[domain.BranchName]runstate.BranchUpdate{},
				BranchesAdded: map[domain.BranchName]domain.SHA{
					domain.NewBranchName("branch-1"): domain.NewSHA("111111"),
					domain.NewBranchName("branch-3"): domain.NewSHA("333333"),
				},
				BranchesRemoved: map[domain.BranchName]domain.SHA{},
				PartialDiff:     runstate.NewPartialDiff(),
			}
			fmt.Printf("WANT: %#v\n", want)
			fmt.Printf("HAVE: %#v\n", have)
			assert.Equal(t, want, have)
		})
		t.Run("branches removed", func(t *testing.T) {
			t.Parallel()
			before := runstate.Snapshot{
				Branches: domain.BranchInfos{
					domain.BranchInfo{
						Name:       domain.NewLocalBranchName("branch-1"),
						InitialSHA: domain.NewSHA("111111"),
					},
					domain.BranchInfo{
						Name:       domain.NewLocalBranchName("branch-2"),
						InitialSHA: domain.NewSHA("222222"),
					},
					domain.BranchInfo{
						Name:       domain.NewLocalBranchName("branch-3"),
						InitialSHA: domain.NewSHA("333333"),
					},
				},
			}
			after := runstate.Snapshot{
				Branches: domain.BranchInfos{
					domain.BranchInfo{
						Name:       domain.NewLocalBranchName("branch-2"),
						InitialSHA: domain.NewSHA("222222"),
					},
				},
			}
			have := after.Diff(before)
			want := runstate.Diff{
				BranchesUpdated: map[domain.BranchName]runstate.BranchUpdate{},
				BranchesAdded:   map[domain.BranchName]domain.SHA{},
				BranchesRemoved: map[domain.BranchName]domain.SHA{
					domain.NewBranchName("branch-1"): domain.NewSHA("111111"),
					domain.NewBranchName("branch-3"): domain.NewSHA("333333"),
				},
				PartialDiff: runstate.NewPartialDiff(),
			}
			fmt.Printf("WANT: %#v\n", want)
			fmt.Printf("HAVE: %#v\n", have)
			assert.Equal(t, want, have)
		})
		t.Run("branches updated", func(t *testing.T) {
			t.Parallel()
			before := runstate.Snapshot{
				Branches: domain.BranchInfos{
					domain.BranchInfo{
						Name:       domain.NewLocalBranchName("branch-1"),
						InitialSHA: domain.NewSHA("111111"),
					},
					domain.BranchInfo{
						Name:       domain.NewLocalBranchName("branch-2"),
						InitialSHA: domain.NewSHA("222222"),
					},
					domain.BranchInfo{
						Name:       domain.NewLocalBranchName("branch-3"),
						InitialSHA: domain.NewSHA("333333"),
					},
				},
			}
			after := runstate.Snapshot{
				Branches: domain.BranchInfos{
					domain.BranchInfo{
						Name:       domain.NewLocalBranchName("branch-1"),
						InitialSHA: domain.NewSHA("111111"),
					},
					domain.BranchInfo{
						Name:       domain.NewLocalBranchName("branch-2"),
						InitialSHA: domain.NewSHA("444444"),
					},
					domain.BranchInfo{
						Name:       domain.NewLocalBranchName("branch-3"),
						InitialSHA: domain.NewSHA("333333"),
					},
				},
			}
			have := after.Diff(before)
			want := runstate.Diff{
				BranchesUpdated: map[domain.BranchName]runstate.BranchUpdate{
					domain.NewBranchName("branch-2"): runstate.BranchUpdate{
						OriginalSHA: domain.NewSHA("222222"),
						FinalSHA:    domain.NewSHA("444444"),
					},
				},
				BranchesAdded:   map[domain.BranchName]domain.SHA{},
				BranchesRemoved: map[domain.BranchName]domain.SHA{},
				PartialDiff:     runstate.NewPartialDiff(),
			}
			fmt.Printf("WANT: %#v\n", want)
			fmt.Printf("HAVE: %#v\n", have)
			assert.Equal(t, want, have)
		})
		t.Run("config added", func(t *testing.T) {
			t.Parallel()
		})
		t.Run("config removed", func(t *testing.T) {
			t.Parallel()
		})
		t.Run("config updated", func(t *testing.T) {
			t.Parallel()
		})
		t.Run("complex example", func(t *testing.T) {
			t.Parallel()
		})
	})
}
