package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v12/src/cli/dialog/components"
	"github.com/git-town/git-town/v12/src/config/configdomain"
	"github.com/git-town/git-town/v12/src/messages"
)

const (
	giteaTokenTitle = `Gitea API token`
	giteaTokenHelp  = `
If you have an API token for Gitea,
and want to ship branches from the CLI,
please enter it now.

It's okay to leave this empty.

`
)

// GiteaToken lets the user enter the Gitea API token.
func GiteaToken(oldValue configdomain.GiteaToken, inputs components.TestInput) (configdomain.GiteaToken, bool, error) {
	token, aborted, err := components.TextField(components.TextFieldArgs{
		ExistingValue: oldValue.String(),
		Help:          giteaTokenHelp,
		Prompt:        "Your Gitea API token: ",
		TestInput:     inputs,
		Title:         giteaTokenTitle,
	})
	fmt.Printf(messages.GiteaToken, components.FormattedSecret(token, aborted))
	return configdomain.GiteaToken(token), aborted, err
}
