package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v16/internal/cli/dialog/components"
	"github.com/git-town/git-town/v16/internal/cli/dialog/components/list"
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/gohacks/slice"
	"github.com/git-town/git-town/v16/internal/messages"
)

const (
	defaultBranchTypeTitle = `Default branch type`
	DefaultBranchTypeHelp  = `
Which type should Git Town assume for branches whose type isn't specified?

When changing this, you should also set the "feature-regex" setting.

`
)

func DefaultBranchType(existingValue configdomain.DefaultBranchType, inputs components.TestInput) (configdomain.DefaultBranchType, bool, error) {
	options := []configdomain.BranchType{
		configdomain.BranchTypeContributionBranch,
		configdomain.BranchTypeFeatureBranch,
		configdomain.BranchTypeObservedBranch,
		configdomain.BranchTypeParkedBranch,
		configdomain.BranchTypePrototypeBranch,
	}
	cursor := slice.Index(options, existingValue.BranchType).GetOrElse(0)
	selection, aborted, err := components.RadioList(list.NewEntries(options...), cursor, defaultBranchTypeTitle, DefaultBranchTypeHelp, inputs)
	fmt.Printf(messages.DefaultBranchType, components.FormattedSelection(selection.String(), aborted))
	return configdomain.DefaultBranchType{BranchType: selection}, aborted, err
}
