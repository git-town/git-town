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
				Branches: map[domain.BranchName]domain.SHA{
					domain.NewBranchName("branch-1"): domain.NewSHA("111111"),
				},
				Config: map[string]string{},
			}
			after := runstate.Snapshot{
				Branches: map[domain.BranchName]domain.SHA{
					domain.NewBranchName("branch-1"): domain.NewSHA("111111"),
					domain.NewBranchName("branch-2"): domain.NewSHA("222222"),
				},
				Config: map[string]string{},
			}
			have := after.Diff(before)
			want := runstate.Diff{
				BranchesUpdated: map[domain.BranchName]runstate.BranchUpdate{},
				BranchesAdded: map[domain.BranchName]domain.SHA{
					domain.NewBranchName("branch-2"): domain.NewSHA("222222"),
				},
				BranchesRemoved: map[domain.BranchName]domain.SHA{},
				ConfigUpdated:   map[string]runstate.ConfigUpdate{},
				ConfigAdded:     map[string]string{},
				ConfigRemoved:   map[string]string{},
			}
			fmt.Printf("WANT: %#v\n", want)
			fmt.Printf("HAVE: %#v\n", have)
			assert.Equal(t, want, have)
		})
		t.Run("branches removed", func(t *testing.T) {
			t.Parallel()
		})
		t.Run("branches updated", func(t *testing.T) {
			t.Parallel()
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
