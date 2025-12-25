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
	bitbucketUsernameTitle = `Bitbucket username`
	bitbucketUsernameHelp  = `
Git Town can update pull requests
and ship branches on Bitbucket for you.
To enable this,
please enter your Bitbucket username.

If you leave this empty,
Git Town will not use the Bitbucket API.

`
)

func BitbucketUsername(args Args[forgedomain.BitbucketUsername]) (Option[forgedomain.BitbucketUsername], dialogdomain.Exit, error) {
	input, exit, err := dialogcomponents.TextField(dialogcomponents.TextFieldArgs{
		DialogName:    "bitbucket-username",
		ExistingValue: args.Local.Or(args.Global).StringOr(""),
		Help:          bitbucketUsernameHelp,
		Inputs:        args.Inputs,
		Prompt:        messages.BitBucketUsernamePrompt,
		Title:         bitbucketUsernameTitle,
	})
	newValue := forgedomain.ParseBitbucketUsername(input)
	if args.Global.Equal(newValue) {
		// the user has entered the global value --> keep using the global value, don't store the local value
		newValue = None[forgedomain.BitbucketUsername]()
	}
	fmt.Printf(messages.BitBucketUsernameResult, dialogcomponents.FormattedOption(newValue, args.Global.IsSome(), exit))
	return newValue, exit, err
}
