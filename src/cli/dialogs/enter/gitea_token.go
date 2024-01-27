package enter

import (
	"fmt"

	"github.com/git-town/git-town/v11/src/cli/dialogs/dialog"
	"github.com/git-town/git-town/v11/src/config/configdomain"
)

const enterGiteaTokenHelp = `
If you have an API token for Gitea,
and want to ship branches from the CLI,
please enter it now.

It's okay to leave this empty.

`

// GiteaToken lets the user enter the Gitea API token.
func GiteaToken(oldValue configdomain.GiteaToken, inputs dialog.TestInput) (configdomain.GiteaToken, bool, error) {
	token, aborted, err := dialog.TextField(dialog.TextFieldArgs{
		ExistingValue: oldValue.String(),
		Help:          enterGiteaTokenHelp,
		Prompt:        "Your Gitea API token: ",
		TestInput:     inputs,
	})
	fmt.Printf("Gitea token: %s\n", dialog.FormattedToken(token, aborted))
	return configdomain.GiteaToken(token), aborted, err
}
