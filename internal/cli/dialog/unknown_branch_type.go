package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents/list"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v22/internal/config"
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/messages"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

const (
	unknownBranchTypeTitle = `Unknown branch type`
	UnknownBranchTypeHelp  = `
Select the type branches get if
Git Town cannot determine their type any other way.

`
)

func UnknownBranchType(args Args[configdomain.UnknownBranchType]) (Option[configdomain.UnknownBranchType], dialogdomain.Exit, error) {
	entries := make(list.Entries[Option[configdomain.UnknownBranchType]], 0, 5)
	globalValue, hasGlobal := args.Global.Get()
	if hasGlobal {
		entries = append(entries, list.Entry[Option[configdomain.UnknownBranchType]]{
			Data: None[configdomain.UnknownBranchType](),
			Text: fmt.Sprintf("use global setting (%s)", globalValue),
		})
	}
	entries = appendEntry(entries, configdomain.BranchTypeContributionBranch)
	entries = appendEntry(entries, configdomain.BranchTypeFeatureBranch)
	entries = appendEntry(entries, configdomain.BranchTypeObservedBranch)
	entries = appendEntry(entries, configdomain.BranchTypeParkedBranch)
	entries = appendEntry(entries, configdomain.BranchTypePrototypeBranch)
	var cursor int
	local, hasLocal := args.Local.Get()
	switch {
	case hasLocal:
		cursor = entries.IndexOf(Some(local))
	case hasGlobal:
		cursor = 0
	default:
		// neither local nor global --> preselect the default value
		cursor = entries.IndexOf(Some(config.DefaultNormalConfig().UnknownBranchType))
	}
	selection, exit, err := dialogcomponents.RadioList(entries, cursor, unknownBranchTypeTitle, UnknownBranchTypeHelp, args.Inputs, "unknown-branch-type")
	fmt.Printf(messages.UnknownBranchType, dialogcomponents.FormattedOption(selection, args.Global.IsSome(), exit))
	return selection, exit, err
}

func appendEntry(entries list.Entries[Option[configdomain.UnknownBranchType]], branchType configdomain.BranchType) list.Entries[Option[configdomain.UnknownBranchType]] {
	return append(entries, list.Entry[Option[configdomain.UnknownBranchType]]{
		Data: Some(configdomain.UnknownBranchType(branchType)),
		Text: branchType.String(),
	})
}
