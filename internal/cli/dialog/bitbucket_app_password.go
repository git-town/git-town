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
	bitbucketAppPasswordTitle = `Bitbucket App Password`
	bitbucketAppPasswordHelp  = `
Git Town can update pull requests
and ship branches on Bitbucket for you.
To enable this, please enter
a Bitbucket App Password.
This is not your normal account password.
More info at
https://www.git-town.com/preferences/bitbucket-app-password.

If you leave this empty,
Git Town will not use the Bitbucket API.

`
)

// BitbucketAppPassword lets the user enter the Bitbucket API token.
func BitbucketAppPassword(args Args[forgedomain.BitbucketAppPassword]) (Option[forgedomain.BitbucketAppPassword], dialogdomain.Exit, error) {
	input, exit, err := dialogcomponents.TextField(dialogcomponents.TextFieldArgs{
		DialogName:    "bitbucket-app-password",
		ExistingValue: args.Local.Or(args.Global).String(),
		Help:          bitbucketAppPasswordHelp,
		Inputs:        args.Inputs,
		Prompt:        messages.BitBucketAppPasswordPrompt,
		Title:         bitbucketAppPasswordTitle,
	})
	newValue := forgedomain.ParseBitbucketAppPassword(input)
	if args.Global.Equal(newValue) {
		// the user has entered the global value --> keep using the global value, don't store the local value
		newValue = None[forgedomain.BitbucketAppPassword]()
	}
	fmt.Printf(messages.BitBucketAppPasswordResult, dialogcomponents.FormattedOption(newValue, args.Global.IsSome(), exit))
	return newValue, exit, err
}
