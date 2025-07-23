package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents/list"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v21/internal/messages"
)

const (
	credentialsNoAccessTitle = "Credentials do not grant API access"
	credentialsNoAccessHelp  = `
The credentials you have entered
seem to not allow API access.

API error message: %v

`
)

func CredentialsNoAccess(connectorError error, inputs dialogcomponents.Inputs) (repeat bool, exit dialogdomain.Exit, err error) {
	entries := list.NewEntries(
		CredentialsNoAccessChoiceRetry,
		CredentialsNoAccessChoiceIgnore,
	)
	selection, exit, err := dialogcomponents.RadioList(entries, 0, credentialsNoAccessTitle, fmt.Sprintf(credentialsNoAccessHelp, connectorError), inputs, "credentials-no-access-to-api")
	fmt.Printf(messages.CredentialsNoAccess, dialogcomponents.FormattedSelection(selection.String(), exit))
	return selection.Repeat(), exit, err
}

type CredentialsNoAccessChoice string

const (
	CredentialsNoAccessChoiceRetry  CredentialsNoAccessChoice = "enter the credentials again"
	CredentialsNoAccessChoiceIgnore CredentialsNoAccessChoice = "store these credentials and continue"
)

func (self CredentialsNoAccessChoice) Repeat() bool {
	switch self {
	case CredentialsNoAccessChoiceRetry:
		return true
	case CredentialsNoAccessChoiceIgnore:
		return false
	}
	panic("unhandled choice")
}

func (self CredentialsNoAccessChoice) String() string {
	return string(self)
}
