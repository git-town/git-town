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
	proposalBreadcrumbSingleTitle = `Show lineage for single stacks`
	proposalBreadcrumbSingleHelp  = `
Should Git Town display breadcrumbs for stacks that contain only a single branch?

`
)

func ProposalBreadcrumbSingle(args Args[forgedomain.ProposalBreadcrumbSingle]) (Option[forgedomain.ProposalBreadcrumbSingle], dialogdomain.Exit, error) {
	entries := list.Entries[Option[forgedomain.ProposalBreadcrumbSingle]]{}
	if global, hasGlobal := args.Global.Get(); hasGlobal {
		entries = append(entries, list.Entry[Option[forgedomain.ProposalBreadcrumbSingle]]{
			Data: None[forgedomain.ProposalBreadcrumbSingle](),
			Text: fmt.Sprintf(messages.DialogUseGlobalValue, global),
		})
	}
	entries = append(entries, list.Entries[Option[forgedomain.ProposalBreadcrumbSingle]]{
		{
			Data: Some(forgedomain.ProposalBreadcrumbSingle(true)),
			Text: "yes, embed breadcrumbs also for single branches",
		},
		{
			Data: Some(forgedomain.ProposalBreadcrumbSingle(false)),
			Text: "no, only embed breadcrums if the stack has at least 2 branches",
		},
	}...)
	defaultPos := entries.IndexOf(args.Local)
	selection, exit, err := dialogcomponents.RadioList(entries, defaultPos, proposalBreadcrumbSingleTitle, proposalBreadcrumbSingleHelp, args.Inputs, "proposal-breadcrumb-single")
	fmt.Printf(messages.ProposalBreadcrumbSingle, dialogcomponents.FormattedOption(selection, args.Global.IsSome(), exit))
	return selection, exit, err
}
