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
	proposalsBreadcrumbSingleTitle = `Show lineage for single stacks`
	proposalBreadcrumbSingleHelp   = `
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
			Text: "yes, embed lineage for single stacks into proposals",
		},
		{
			Data: Some(forgedomain.ProposalBreadcrumbSingle(false)),
			Text: "no, don't embed lineage for single stacks into proposals",
		},
	}...)
	defaultPos := entries.IndexOf(args.Local)
	selection, exit, err := dialogcomponents.RadioList(entries, defaultPos, proposalsBreadcrumbSingleTitle, proposalBreadcrumbSingleHelp, args.Inputs, "proposals-show-lineage-single-stack")
	fmt.Printf(messages.ProposalBreadcrumbSingleStack, dialogcomponents.FormattedOption(selection, args.Global.IsSome(), exit))
	return selection, exit, err
}
