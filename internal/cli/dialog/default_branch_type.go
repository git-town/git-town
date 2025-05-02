package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v20/internal/cli/dialog/components"
	"github.com/git-town/git-town/v20/internal/cli/dialog/components/list"
	"github.com/git-town/git-town/v20/internal/config/configdomain"
	"github.com/git-town/git-town/v20/internal/gohacks/slice"
	"github.com/git-town/git-town/v20/internal/messages"
)

const (
	defaultBranchTypeTitle = `Default branch type`
	DefaultBranchTypeHelp  = `
Select the type Git Town should assume for new branches.

If you change this, consider also configuring
the "feature-regex" setting on the next screen.

`
)

func DefaultBranchType(existingValue configdomain.BranchType, inputs components.TestInput) (configdomain.BranchType, bool, error) {
	options := []configdomain.BranchType{
		configdomain.BranchTypeContributionBranch,
		configdomain.BranchTypeFeatureBranch,
		configdomain.BranchTypeObservedBranch,
		configdomain.BranchTypeParkedBranch,
		configdomain.BranchTypePrototypeBranch,
	}
	cursor := slice.Index(options, existingValue).GetOrElse(0)
	selection, aborted, err := components.RadioList(list.NewEntries(options...), cursor, defaultBranchTypeTitle, DefaultBranchTypeHelp, inputs)
	fmt.Printf(messages.DefaultBranchType, components.FormattedSelection(selection.String(), aborted))
	return selection, aborted, err
}
