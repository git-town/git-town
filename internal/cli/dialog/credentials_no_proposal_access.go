package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents/list"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v22/internal/messages"
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

func CredentialsNoProposalAccess(connectorError error, inputs dialogcomponents.Inputs) (repeat bool, exit dialogdomain.Exit, err error) {
	entries := list.NewEntries(
		CredentialsNoAccessChoiceRetry,
		CredentialsNoAccessChoiceIgnore,
	)
	selection, exit, err := dialogcomponents.RadioList(entries, 0, credentialsNoProposalAccessTitle, fmt.Sprintf(credentialsNoProposalAccessHelp, connectorError), inputs, "credentials-no-access-to-proposal")
	fmt.Printf(messages.CredentialsNoAccess, dialogcomponents.FormattedSelection(selection.String(), exit))
	return selection.Repeat(), exit, err
}
