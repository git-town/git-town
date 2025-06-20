package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/cli/dialog/components"
	"github.com/git-town/git-town/v21/internal/cli/dialog/components/list"
	"github.com/git-town/git-town/v21/internal/messages"
)

const (
	credentialsNoProposalAccessTitle = "Credentials do not grant access to proposals"
	credentialsNoProposalAccessHelp  = `
The credentials you have entered allow API access, but not to access proposals.

Received error: %s
`
)

func CredentialsNoProposalAccess(connectorError error, inputs components.TestInput) (repeat bool, aborted bool, err error) {
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
	selection, aborted, err := components.RadioList(entries, defaultPos, credentialsNoProposalAccessTitle, fmt.Sprintf(credentialsNoProposalAccessHelp, connectorError), inputs)
	if err != nil || aborted {
		return selection.Repeat(), aborted, err
	}
	fmt.Printf(messages.CredentialsNoAccess, components.FormattedSelection(selection.String(), aborted))
	return selection.Repeat(), aborted, err
}
