package datatable_test

import (
	"testing"

	"github.com/cucumber/godog"
	messages "github.com/cucumber/messages/go/v21"
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/git-town/git-town/v14/test/datatable"
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
					},
				},
				{
					Cells: []*messages.PickleTableCell{
						{Value: "main"},
						{Value: "main"},
						{Value: ""},
					},
				},
				{
					Cells: []*messages.PickleTableCell{
						{Value: "feature-1"},
						{Value: "feature"},
						{Value: "main"},
					},
				},
			},
		}
		have := datatable.ParseBranchSetupTable(give)
		want := []datatable.BranchSetup{
			{
				Name:       "main",
				BranchType: configdomain.BranchTypeMainBranch,
				Parent:     None[gitdomain.LocalBranchName](),
			},
			{
				Name:       "feature-1",
				BranchType: configdomain.BranchTypeFeatureBranch,
				Parent:     Some(gitdomain.NewLocalBranchName("main")),
			},
		}
		must.Eq(t, want, have)
	})
}
