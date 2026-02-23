package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents/list"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/messages"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

const (
	proposalBreadcrumbDirectionTitle = `Proposal Breadcrumb direction`
	proposalBreadcrumbDirectionHelp  = `
	Which direction should proposal breadcrumbs have?

	See https://www.git-town.com/how-to/proposal-breadcrumb-direction.html for details.
`
)

func ProposalBreadcrumbDirection(args Args[configdomain.ProposalBreadcrumbDirection]) (Option[configdomain.ProposalBreadcrumbDirection], dialogdomain.Exit, error) {
	entries := list.Entries[Option[configdomain.ProposalBreadcrumbDirection]]{}
	if global, hasGlobal := args.Global.Get(); hasGlobal {
		entries = append(entries, list.Entry[Option[configdomain.ProposalBreadcrumbDirection]]{
			Data: None[configdomain.ProposalBreadcrumbDirection](),
			Text: fmt.Sprintf(messages.DialogUseGlobalValue, global),
		})
	}

	entries = append(entries, list.Entries[Option[configdomain.ProposalBreadcrumbDirection]]{
		{
			Data: Some(configdomain.ProposalBreadcrumbDirectionDown),
			Text: "down from the root (default)",
		},
		{
			Data: Some(configdomain.ProposalBreadcrumbDirectionUp),
			Text: "up from the root",
		},
	}...)
	defaultPos := entries.IndexOf(args.Local)
	selection, exit, err := dialogcomponents.RadioList(entries, defaultPos, proposalBreadcrumbDirectionTitle, proposalBreadcrumbDirectionHelp, args.Inputs, "proposal-breadcrumb-direction")
	fmt.Printf(messages.ProposalBreadcrumbDirection, dialogcomponents.FormattedOption(selection, args.Global.IsSome(), exit))
	return selection, exit, err
}
