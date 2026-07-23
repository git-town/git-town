package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v24/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v24/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v24/internal/forge/forgedomain"
	"github.com/git-town/git-town/v24/internal/messages"
	. "github.com/git-town/git-town/v24/pkg/prelude"
)

const (
	bitbucketUsernameTitle     = `Bitbucket username`
	bitbucketUsernameHelpCloud = `
Git Town can update pull requests
and ship branches on Bitbucket for you.
To enable this,
please enter your Bitbucket username.
When authenticating with a Bitbucket Cloud
API token, this is the email address
of your Atlassian account.
More info at
https://www.git-town.com/preferences/bitbucket-api-token.

If you leave this empty,
Git Town will not use the Bitbucket API.

`
	bitbucketUsernameHelpDataCenter = `
Git Town can update pull requests
and ship branches on Bitbucket for you.
To enable this,
please enter your Bitbucket username.
More info at
https://www.git-town.com/preferences/bitbucket-api-token.

If you leave this empty,
Git Town will not use the Bitbucket API.

`
)

func BitbucketUsername(forgeType forgedomain.ForgeType, args Args[forgedomain.BitbucketUsername]) (Option[forgedomain.BitbucketUsername], dialogdomain.Exit, error) {
	input, exit, err := dialogcomponents.TextField(dialogcomponents.TextFieldArgs{
		DialogName:    "bitbucket-username",
		ExistingValue: args.Local.Or(args.Global).StringOr(""),
		Help:          BitbucketUsernameHelp(forgeType),
		Inputs:        args.Inputs,
		Interactive:   args.Interactive,
		Prompt:        messages.BitbucketUsernamePrompt,
		Title:         bitbucketUsernameTitle,
	})
	newValue := forgedomain.ParseBitbucketUsername(input)
	if args.Global.Equal(newValue) {
		// the user has entered the global value --> keep using the global value, don't store the local value
		newValue = None[forgedomain.BitbucketUsername]()
	}
	fmt.Printf(messages.BitbucketUsernameResult, dialogcomponents.FormattedOption(newValue, args.Global.IsSome(), exit))
	return newValue, exit, err
}

// BitbucketUsernameHelp provides the help text for the Bitbucket username dialog
// that matches the given Bitbucket variant.
func BitbucketUsernameHelp(forgeType forgedomain.ForgeType) string {
	if forgeType == forgedomain.ForgeTypeBitbucketDatacenter {
		return bitbucketUsernameHelpDataCenter
	}
	return bitbucketUsernameHelpCloud
}
