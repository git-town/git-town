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
	proposalBreadcrumbTitle = `Proposal Breadcrumb`
	proposalBreadcrumbHelp  = `
	Should proposals contain a breadcrumb of proposals for all branches in the stack?

	See https://www.git-town.com/how-to/proposal-breadcrumb.html for details.
`
)

func ProposalBreadcrumb(args Args[forgedomain.ProposalBreadcrumb]) (Option[forgedomain.ProposalBreadcrumb], dialogdomain.Exit, error) {
	entries := list.Entries[Option[forgedomain.ProposalBreadcrumb]]{}
	if global, hasGlobal := args.Global.Get(); hasGlobal {
		entries = append(entries, list.Entry[Option[forgedomain.ProposalBreadcrumb]]{
			Data: None[forgedomain.ProposalBreadcrumb](),
			Text: fmt.Sprintf(messages.DialogUseGlobalValue, global),
		})
	}

	entries = append(entries, list.Entries[Option[forgedomain.ProposalBreadcrumb]]{
		{
			Data: Some(forgedomain.ProposalBreadcrumbNone),
			Text: "no breadcrumb in proposals, or use the Git Town GitHub Action",
		},
		{
			Data: Some(forgedomain.ProposalBreadcrumbStacks),
			Text: "Git Town CLI embeds the breadcrumbs for stacks containing more than 2 branches into proposals",
		},
		{
			Data: Some(forgedomain.ProposalBreadcrumbBranches),
			Text: "Git Town CLI embeds the breadcrumbs into all proposals",
		},
	}...)
	defaultPos := entries.IndexOf(args.Local)
	selection, exit, err := dialogcomponents.RadioList(entries, defaultPos, proposalBreadcrumbTitle, proposalBreadcrumbHelp, args.Inputs, "proposal-breadcrumb")
	fmt.Printf(messages.ProposalBreadcrumb, dialogcomponents.FormattedOption(selection, args.Global.IsSome(), exit))
	return selection, exit, err
}
