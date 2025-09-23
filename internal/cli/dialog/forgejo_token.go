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
	forgejoTokenTitle = `Forgejo API token`
	forgejoTokenHelp  = `
Git Town can update pull requests
and ship branches on Forgejo-based forges for you.
To enable this, please enter a Forgejo API token.
More info at
https://docs.codeberg.org/advanced/access-token.

If you leave this empty,
Git Town will not use the Forgejo API.

`
)

func ForgejoToken(args Args[forgedomain.ForgejoToken]) (Option[forgedomain.ForgejoToken], dialogdomain.Exit, error) {
	input, exit, err := dialogcomponents.TextField(dialogcomponents.TextFieldArgs{
		DialogName:    "forgejo-token",
		ExistingValue: args.Local.Or(args.Global).StringOr(""),
		Help:          forgejoTokenHelp,
		Inputs:        args.Inputs,
		Prompt:        messages.ForgejoTokenPrompt,
		Title:         forgejoTokenTitle,
	})
	newValue := forgedomain.ParseForgejoToken(input)
	if args.Global.Equal(newValue) {
		// the user has entered the global value --> keep using the global value, don't store the local value
		newValue = None[forgedomain.ForgejoToken]()
	}
	fmt.Printf(messages.ForgejoTokenResult, dialogcomponents.FormattedOption(newValue, args.Global.IsSome(), exit))
	return newValue, exit, err
}
