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
	forgejoTokenTitle = `Forgejo API token`
	forgejoTokenHelp  = `
Git Town can update pull requests
and ship branches on Forgejo-based forges for you.
To enable this, please enter a codeberg API token.
More info at
https://docs.codeberg.org/advanced/access-token.

If you leave this empty,
Git Town will not use the codeberg API.

`
)

func ForgejoToken(args Args[forgedomain.CodebergToken]) (Option[forgedomain.CodebergToken], dialogdomain.Exit, error) {
	input, exit, err := dialogcomponents.TextField(dialogcomponents.TextFieldArgs{
		DialogName:    "codeberg-token",
		ExistingValue: args.Local.Or(args.Global).String(),
		Help:          forgejoTokenHelp,
		Inputs:        args.Inputs,
		Prompt:        messages.CodebergTokenPrompt,
		Title:         forgejoTokenTitle,
	})
	newValue := forgedomain.ParseCodebergToken(input)
	if args.Global.Equal(newValue) {
		// the user has entered the global value --> keep using the global value, don't store the local value
		newValue = None[forgedomain.CodebergToken]()
	}
	fmt.Printf(messages.CodebergTokenResult, dialogcomponents.FormattedOption(newValue, args.Global.IsSome(), exit))
	return newValue, exit, err
}
