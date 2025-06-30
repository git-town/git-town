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
The credentials you have entered
allow API access,
but don't allow access to proposals.

Received error: %v
`
)

func CredentialsNoProposalAccess(connectorError error, inputs components.TestInput) (repeat bool, exit dialogdomain.Exit, err error) {
	entries := list.NewEntries(
		CredentialsNoAccessChoiceRetry,
		CredentialsNoAccessChoiceIgnore,
	)
	selection, exit, err := components.RadioList(entries, 0, credentialsNoProposalAccessTitle, fmt.Sprintf(credentialsNoProposalAccessHelp, connectorError), inputs)
	fmt.Printf(messages.CredentialsNoAccess, components.FormattedSelection(selection.String(), exit))
	return selection.Repeat(), exit, err
}
