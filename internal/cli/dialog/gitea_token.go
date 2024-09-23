package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v16/internal/cli/dialog/components"
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/messages"
	. "github.com/git-town/git-town/v16/pkg/prelude"
)

const (
	giteaTokenTitle = `Gitea API token`
	giteaTokenHelp  = `
Git Town can update pull requests and ship branches on gitea for you.
To enable this, please enter a gitea API token.
More info at https://www.git-town.com/preferences/gitea-token.

If you leave this empty, Git Town will not use the gitea API.

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
