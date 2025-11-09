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
	proposalsShowLineageTitle = `Proposals Show Lineage`
	proposalsShowLineageHelp  = `
How should Git Town share stack lineage in proposals?

Possible Options:
	- none: (default) Git Town should not share stack 
					lineage in proposals
	- cli: Git Town shares or updates stack lineage in proposals
				anytime your stack lineage changes via the cli
	- ci: Git Town manages and shows proposal lineage through ci 
	      integrations, with tools such as https://github.com/git-town/action
`
)

func ProposalsShowLineage(args Args[forgedomain.ProposalsShowLineage]) (Option[forgedomain.ProposalsShowLineage], dialogdomain.Exit, error) {
	entries := list.Entries[Option[forgedomain.ProposalsShowLineage]]{}
	if global, hasGlobal := args.Global.Get(); hasGlobal {
		entries = append(entries, list.Entry[Option[forgedomain.ProposalsShowLineage]]{
			Data: None[forgedomain.ProposalsShowLineage](),
			Text: fmt.Sprintf(messages.DialogUseGlobalValue, global),
		})
	}

	entries = append(entries, list.Entries[Option[forgedomain.ProposalsShowLineage]]{
		{
			Data: Some(forgedomain.ProposalsShowLineageNone),
			Text: "do not show stack lineage in proposals",
		},
		{
			Data: Some(forgedomain.ProposalsShowLineageCLI),
			Text: "cli command execution show the latest stack lineage in a proposal(s)",
		},
		{
			Data: Some(forgedomain.ProposalsShowLineageCI),
			Text: "ci process shows the latest  stack lineage in a proposal(s)",
		},
	}...)
	defaultPos := entries.IndexOf(args.Local)
	selection, exit, err := dialogcomponents.RadioList(entries, defaultPos, proposalsShowLineageTitle, proposalsShowLineageHelp, args.Inputs, "proposals-show-lineage")
	fmt.Printf(messages.ProposalsLineage, dialogcomponents.FormattedOption(selection, args.Global.IsSome(), exit))
	return selection, exit, err
}
