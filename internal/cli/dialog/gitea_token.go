package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v15/internal/cli/dialog/components"
	"github.com/git-town/git-town/v15/internal/config/configdomain"
	. "github.com/git-town/git-town/v15/internal/gohacks/prelude"
	"github.com/git-town/git-town/v15/internal/messages"
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
func GiteaToken(oldValue Option[configdomain.GiteaToken], inputs components.TestInput) (Option[configdomain.GiteaToken], bool, error) {
	text, aborted, err := components.TextField(components.TextFieldArgs{
		ExistingValue: oldValue.String(),
		Help:          giteaTokenHelp,
		Prompt:        "Your Gitea API token: ",
		TestInput:     inputs,
		Title:         giteaTokenTitle,
	})
	fmt.Printf(messages.GiteaToken, components.FormattedSecret(text, aborted))
	return configdomain.ParseGiteaToken(text), aborted, err
}
