package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/cli/dialog/components"
	"github.com/git-town/git-town/v21/internal/cli/dialog/components/list"
	"github.com/git-town/git-town/v21/internal/messages"
)

const (
	credentialsNoAccessTitle = "Credentials do not grant API access"
	credentialsNoAccessHelp  = `
The credentials you have entered seem to not allow API access at all.

Received error: %s
`
)

func CredentialsNoAccess(err error, inputs components.TestInput) (CredentialsNoAccessChoice, bool, error) {
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
	selection, aborted, err := components.RadioList(entries, defaultPos, credentialsNoAccessTitle, credentialsNoAccessHelp, inputs)
	if err != nil || aborted {
		return selection, aborted, err
	}
	fmt.Printf(messages.CredentialsNoAccess, components.FormattedSelection(selection.String(), aborted))
	return selection, aborted, err
}

type CredentialsNoAccessChoice string

const (
	CredentialsNoAccessChoiceRetry  = "retry"
	CredentialsNoAccessChoiceIgnore = "ignore"
)

func (self CredentialsNoAccessChoice) String() string {
	return string(self)
}
