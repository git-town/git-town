package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/cli/dialog/components"
	"github.com/git-town/git-town/v21/internal/cli/dialog/components/list"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/gohacks/slice"
	"github.com/git-town/git-town/v21/internal/messages"
)

const (
	unknownBranchTypeTitle = `Unknown branch type`
	UnknownBranchTypeHelp  = `
Select the type branches get if
Git Town cannot determine their type any other way.

If you set this to something other than "feature",
consider also configuring the "feature-regex" setting
on the next screen.

`
)

func UnknownBranchType(existingValue configdomain.BranchType, inputs components.TestInput) (configdomain.BranchType, dialogdomain.Exit, error) {
	options := []configdomain.BranchType{
		configdomain.BranchTypeContributionBranch,
		configdomain.BranchTypeFeatureBranch,
		configdomain.BranchTypeObservedBranch,
		configdomain.BranchTypeParkedBranch,
		configdomain.BranchTypePrototypeBranch,
	}
	cursor := slice.Index(options, existingValue).GetOrElse(0)
	selection, exit, err := components.RadioList(list.NewEntries(options...), cursor, unknownBranchTypeTitle, UnknownBranchTypeHelp, inputs)
	fmt.Printf(messages.UnknownBranchType, components.FormattedSelection(selection.String(), exit))
	return selection, exit, err
}
