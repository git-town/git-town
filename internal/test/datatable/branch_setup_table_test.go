package datatable_test

import (
	"testing"

	"github.com/cucumber/godog"
	messages "github.com/cucumber/messages/go/v21"
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/test/datatable"
	"github.com/git-town/git-town/v22/internal/test/testgit"
	. "github.com/git-town/git-town/v22/pkg/prelude"
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
				BranchType: Some(configdomain.BranchTypeFeatureBranch),
				Parent:     gitdomain.NewLocalBranchNameOption("main"),
				Locations:  []testgit.Location{testgit.LocationLocal, testgit.LocationOrigin},
			},
			{
				Name:       "feature-2",
				BranchType: Some(configdomain.BranchTypeFeatureBranch),
				Parent:     gitdomain.NewLocalBranchNameOption("main"),
				Locations:  []testgit.Location{testgit.LocationLocal, testgit.LocationOrigin},
			},
		}
		must.Eq(t, want, have)
	})

	t.Run("without parents", func(t *testing.T) {
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
				BranchType: Some(configdomain.BranchTypePerennialBranch),
				Parent:     None[gitdomain.LocalBranchName](),
				Locations:  testgit.Locations{testgit.LocationLocal, testgit.LocationOrigin},
			},
		}
		must.Eq(t, want, have)
	})
}
