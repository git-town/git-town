package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/messages"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

const (
	giteaTokenTitle = `Gitea API token`
	giteaTokenHelp  = `
Git Town can update pull requests
and ship branches on gitea for you.
To enable this, please enter a gitea API token.
More info at
https://www.git-town.com/preferences/gitea-token.

If you leave this empty,
Git Town will not use the gitea API.

`
)

// GiteaToken lets the user enter the Gitea API token.
func GiteaToken(oldValue Option[forgedomain.GiteaToken], inputs dialogcomponents.TestInputs) (Option[forgedomain.GiteaToken], dialogdomain.Exit, error) {
	text, exit, err := dialogcomponents.TextField(dialogcomponents.TextFieldArgs{
		ExistingValue: oldValue.String(),
		Help:          giteaTokenHelp,
		Prompt:        "Your Gitea API token: ",
		TestInputs:    inputs,
		Title:         giteaTokenTitle,
	})
	fmt.Printf(messages.GiteaToken, dialogcomponents.FormattedSecret(text, exit))
	return forgedomain.ParseGiteaToken(text), exit, err
}
