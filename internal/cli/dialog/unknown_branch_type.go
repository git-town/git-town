package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents/list"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v21/internal/config"
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

func UnknownBranchType(unvalidatedConfig config.UnvalidatedConfig, inputs dialogcomponents.TestInputs) (Option[configdomain.UnknownBranchType], dialogdomain.Exit, error) {
	entries := make(list.Entries[Option[configdomain.UnknownBranchType]], 0, 5)
	if globalValue, has := unvalidatedConfig.GitGlobal.UnknownBranchType.Get(); has {
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
	defaultPos := determinePos(entries, unvalidatedConfig)
	selection, exit, err := dialogcomponents.RadioList(entries, defaultPos, unknownBranchTypeTitle, UnknownBranchTypeHelp, inputs, "unknown-branch-type")
	value, has := selection.Get()
	fmt.Printf(messages.UnknownBranchType, dialogcomponents.FormattedOptionalSelection(value, has, exit))
	return selection, exit, err
}

func appendEntry(entries list.Entries[Option[configdomain.UnknownBranchType]], branchType configdomain.BranchType) list.Entries[Option[configdomain.UnknownBranchType]] {
	return append(entries, list.Entry[Option[configdomain.UnknownBranchType]]{
		Data: Some(configdomain.UnknownBranchType(branchType)),
		Text: branchType.String(),
	})
}

func determinePos(entries list.Entries[Option[configdomain.UnknownBranchType]], unvalidatedConfig config.UnvalidatedConfig) int {
	if localValue, has := unvalidatedConfig.GitLocal.UnknownBranchType.Get(); has {
		return entries.IndexOf(Some(localValue))
	}
	if unvalidatedConfig.GitGlobal.UnknownBranchType.IsSome() {
		return 0
	}
	return entries.IndexOf(Some(unvalidatedConfig.Defaults.UnknownBranchType))
}
