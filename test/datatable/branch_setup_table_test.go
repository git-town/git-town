package datatable_test

import (
	"testing"

	"github.com/cucumber/godog"
	messages "github.com/cucumber/messages/go/v21"
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/git-town/git-town/v14/test/datatable"
	"github.com/git-town/git-town/v14/test/git"
	"github.com/shoenig/test/must"
)

func TestParseBranchSetupTable(t *testing.T) {
	t.Parallel()

	t.Run("normal table", func(t *testing.T) {
		t.Parallel()
		give := &godog.Table{
			Rows: []*messages.PickleTableRow{
				{
					Cells: []*messages.PickleTableCell{
						{Value: "NAME"},
						{Value: "TYPE"},
						{Value: "PARENT"},
						{Value: "LOCATIONS"},
					},
				},
				{
					Cells: []*messages.PickleTableCell{
						{Value: "feature-1"},
						{Value: "feature"},
						{Value: "main"},
						{Value: "local, origin"},
					},
				},
				{
					Cells: []*messages.PickleTableCell{
						{Value: "feature-2"},
						{Value: "feature"},
						{Value: "main"},
						{Value: ""},
					},
				},
			},
		}
		have := datatable.ParseBranchSetupTable(give)
		want := []datatable.BranchSetup{
			{
				Name:       "feature-1",
				BranchType: configdomain.BranchTypeFeatureBranch,
				Parent:     Some(gitdomain.NewLocalBranchName("main")),
				Locations:  []git.Location{git.LocationLocal, git.LocationOrigin},
			},
			{
				Name:       "feature-2",
				BranchType: configdomain.BranchTypeFeatureBranch,
				Parent:     Some(gitdomain.NewLocalBranchName("main")),
				Locations:  []git.Location{git.LocationLocal, git.LocationOrigin},
			},
		}
		must.Eq(t, want, have)
	})

	t.Run("no parents given", func(t *testing.T) {
		t.Parallel()
		give := &godog.Table{
			Rows: []*messages.PickleTableRow{
				{
					Cells: []*messages.PickleTableCell{
						{Value: "NAME"},
						{Value: "TYPE"},
						{Value: "LOCATIONS"},
					},
				},
				{
					Cells: []*messages.PickleTableCell{
						{Value: "staging"},
						{Value: "perennial"},
						{Value: "local, origin"},
					},
				},
			},
		}
		have := datatable.ParseBranchSetupTable(give)
		want := []datatable.BranchSetup{
			{
				Name:       "staging",
				BranchType: configdomain.BranchTypePerennialBranch,
				Parent:     None[gitdomain.LocalBranchName](),
				Locations:  git.Locations{git.LocationLocal, git.LocationOrigin},
			},
		}
		must.Eq(t, want, have)
	})
}
