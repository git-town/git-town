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
	unknownBranchTypeTitle = `Unknown branch type`
	UnknownBranchTypeHelp  = `
Select the type branches get if
feature-regex, observed-regex, contribution-regex,
perennial-regex, the perennial branch list, don't match.

If you set this to something other than "feature",
consider also configuring the "feature-regex" setting
on the next screen.

`
)

func UnknownBranchType(existingValue configdomain.BranchType, inputs components.TestInput) (configdomain.BranchType, bool, error) {
	options := []configdomain.BranchType{
		configdomain.BranchTypeContributionBranch,
		configdomain.BranchTypeFeatureBranch,
		configdomain.BranchTypeObservedBranch,
		configdomain.BranchTypeParkedBranch,
		configdomain.BranchTypePrototypeBranch,
	}
	cursor := slice.Index(options, existingValue).GetOrElse(0)
	selection, aborted, err := components.RadioList(list.NewEntries(options...), cursor, unknownBranchTypeTitle, UnknownBranchTypeHelp, inputs)
	fmt.Printf(messages.UnknownBranchType, components.FormattedSelection(selection.String(), aborted))
	return selection, aborted, err
}
