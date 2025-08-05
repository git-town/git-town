package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents/list"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/messages"
	. "github.com/git-town/git-town/v21/pkg/prelude"
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
	if globalValue, hasGlobal := args.Global.Get(); hasGlobal {
		entries = append(entries, list.Entry[Option[configdomain.UnknownBranchType]]{
			Data: None[configdomain.UnknownBranchType](),
			Text: fmt.Sprintf(messages.DialogUseGlobalValue, globalValue),
		})
	} else {
		entries = append(entries, list.Entry[Option[configdomain.UnknownBranchType]]{
			Data: None[configdomain.UnknownBranchType](),
			Text: fmt.Sprintf(messages.DialogUseDefaultValue, args.Defaults),
		})
	}
	entries = appendEntry(entries, configdomain.BranchTypeContributionBranch)
	entries = appendEntry(entries, configdomain.BranchTypeFeatureBranch)
	entries = appendEntry(entries, configdomain.BranchTypeObservedBranch)
	entries = appendEntry(entries, configdomain.BranchTypeParkedBranch)
	entries = appendEntry(entries, configdomain.BranchTypePrototypeBranch)
	cursor := 0
	if local, hasLocal := args.Local.Get(); hasLocal {
		cursor = entries.IndexOf(Some(local))
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
