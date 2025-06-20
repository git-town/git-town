package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/cli/dialog/components"
	"github.com/git-town/git-town/v21/internal/cli/dialog/components/list"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v21/internal/messages"
)

const (
	credentialsNoAccessTitle = "Credentials do not grant API access"
	credentialsNoAccessHelp  = `
The credentials you have entered seem to not allow API access at all.

API error message: %v

`
)

func CredentialsNoAccess(connectorError error, inputs components.TestInput) (repeat bool, exit dialogdomain.Exit, err error) {
	entries := list.Entries[CredentialsNoAccessChoice]{
		{
			Data: CredentialsNoAccessChoiceRetry,
			Text: CredentialsNoAccessChoiceRetry,
		},
		{
			Data: CredentialsNoAccessChoiceIgnore,
			Text: CredentialsNoAccessChoiceIgnore,
		},
	}
	defaultPos := entries.IndexOf(CredentialsNoAccessChoiceRetry)
	selection, exit, err := components.RadioList(entries, defaultPos, credentialsNoAccessTitle, fmt.Sprintf(credentialsNoAccessHelp, connectorError), inputs)
	if err != nil || exit {
		return false, exit, err
	}
	fmt.Printf(messages.CredentialsNoAccess, components.FormattedSelection(selection.String(), exit))
	return selection.Repeat(), exit, err
}

type CredentialsNoAccessChoice string

const (
	CredentialsNoAccessChoiceRetry  = "enter the credentials again"
	CredentialsNoAccessChoiceIgnore = "store these credentials and continue"
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
