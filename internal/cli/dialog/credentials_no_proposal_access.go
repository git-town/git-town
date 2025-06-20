package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/cli/dialog/components"
	"github.com/git-town/git-town/v21/internal/cli/dialog/components/list"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v21/internal/messages"
)

const (
	credentialsNoProposalAccessTitle = "Credentials do not grant access to proposals"
	credentialsNoProposalAccessHelp  = `
The credentials you have entered allow API access, but don't allow access to proposals.

Received error: %v
`
)

func CredentialsNoProposalAccess(connectorError error, inputs components.TestInput) (repeat bool, exit dialogdomain.Exit, err error) {
	entries := list.Entries[CredentialsNoAccessChoice]{
		{
			Data: CredentialsNoAccessChoiceRetry,
			Text: `enter the credentials again`,
		},
		{
			Data: CredentialsNoAccessChoiceIgnore,
			Text: `store these credentials and continue`,
		},
	}
	defaultPos := entries.IndexOf(CredentialsNoAccessChoiceRetry)
	selection, exit, err := components.RadioList(entries, defaultPos, credentialsNoProposalAccessTitle, fmt.Sprintf(credentialsNoProposalAccessHelp, connectorError), inputs)
	if err != nil || exit {
		return selection.Repeat(), exit, err
	}
	fmt.Printf(messages.CredentialsNoAccess, components.FormattedSelection(selection.String(), exit))
	return selection.Repeat(), exit, err
}
