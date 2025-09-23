package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/messages"
	. "github.com/git-town/git-town/v22/pkg/prelude"
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

func GiteaToken(args Args[forgedomain.GiteaToken]) (Option[forgedomain.GiteaToken], dialogdomain.Exit, error) {
	input, exit, err := dialogcomponents.TextField(dialogcomponents.TextFieldArgs{
		DialogName:    "gitea-token",
		ExistingValue: args.Local.Or(args.Global).StringOr(""),
		Help:          giteaTokenHelp,
		Inputs:        args.Inputs,
		Prompt:        messages.GiteaTokenPrompt,
		Title:         giteaTokenTitle,
	})
	newValue := forgedomain.ParseGiteaToken(input)
	if args.Global.Equal(newValue) {
		// the user has entered the global value --> keep using the global value, don't store the local value
		newValue = None[forgedomain.GiteaToken]()
	}
	fmt.Printf(messages.GiteaTokenResult, dialogcomponents.FormattedOption(newValue, args.Global.IsSome(), exit))
	return newValue, exit, err
}
