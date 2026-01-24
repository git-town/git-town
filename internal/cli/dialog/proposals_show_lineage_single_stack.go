package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents/list"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/messages"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

const (
	proposalsShowLineageSingleStackTitle = `Show lineage for single stacks`
	proposalShowLineageSingleStackHelp   = `
Should "git town sync" sync Git tags with origin?

`
)

func ProposalShowLineageSingleStack(args Args[forgedomain.ProposalsShowLineageSingleStack]) (Option[forgedomain.ProposalsShowLineageSingleStack], dialogdomain.Exit, error) {
	entries := list.Entries[Option[forgedomain.ProposalsShowLineageSingleStack]]{}
	if global, hasGlobal := args.Global.Get(); hasGlobal {
		entries = append(entries, list.Entry[Option[forgedomain.ProposalsShowLineageSingleStack]]{
			Data: None[forgedomain.ProposalsShowLineageSingleStack](),
			Text: fmt.Sprintf(messages.DialogUseGlobalValue, global),
		})
	}
	entries = append(entries, list.Entries[Option[forgedomain.ProposalsShowLineageSingleStack]]{
		{
			Data: Some(forgedomain.ProposalsShowLineageSingleStack(true)),
			Text: "yes, embed lineage for single stacks into proposals",
		},
		{
			Data: Some(forgedomain.ProposalsShowLineageSingleStack(false)),
			Text: "no, don't embed lineage for single stacks into proposals",
		},
	}...)
	defaultPos := entries.IndexOf(args.Local)
	selection, exit, err := dialogcomponents.RadioList(entries, defaultPos, proposalsShowLineageSingleStackTitle, proposalShowLineageSingleStackHelp, args.Inputs, "proposals-show-lineage-single-stack")
	fmt.Printf(messages.ProposalBreadcrumbSingleStack, dialogcomponents.FormattedOption(selection, args.Global.IsSome(), exit))
	return selection, exit, err
}
