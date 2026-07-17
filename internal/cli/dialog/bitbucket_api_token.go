package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v23/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v23/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v23/internal/forge/forgedomain"
	"github.com/git-town/git-town/v23/internal/messages"
	. "github.com/git-town/git-town/v23/pkg/prelude"
)

const (
	bitbucketAPITokenTitle     = `Bitbucket API token`
	bitbucketAPITokenHelpCloud = `
Git Town can update pull requests
and ship branches on Bitbucket for you.
To enable this, please enter
a Bitbucket API token with scopes.
This is not your normal account password.
You can create one at
https://id.atlassian.com/manage-profile/security/api-tokens.
More info at
https://www.git-town.com/preferences/bitbucket-api-token.

If you leave this empty,
Git Town will not use the Bitbucket API.

`
	bitbucketAPITokenHelpDataCenter = `
Git Town can update pull requests
and ship branches on Bitbucket for you.
To enable this, please enter
an HTTP access token with
"Project read" and "Repository write"
permissions.
This is not your normal account password.
You can create one in your Bitbucket instance
under Profile picture > Manage account >
HTTP access tokens.
More info at
https://www.git-town.com/preferences/bitbucket-api-token.

If you leave this empty,
Git Town will not use the Bitbucket API.

`
)

// BitbucketAPIToken lets the user enter the Bitbucket API token.
func BitbucketAPIToken(forgeType forgedomain.ForgeType, args Args[forgedomain.BitbucketAPIToken]) (Option[forgedomain.BitbucketAPIToken], dialogdomain.Exit, error) {
	input, exit, err := dialogcomponents.TextField(dialogcomponents.TextFieldArgs{
		DialogName:    "bitbucket-api-token",
		ExistingValue: args.Local.Or(args.Global).StringOr(""),
		Help:          BitbucketAPITokenHelp(forgeType),
		Inputs:        args.Inputs,
		Interactive:   args.Interactive,
		Prompt:        messages.BitbucketAPITokenPrompt,
		Title:         bitbucketAPITokenTitle,
	})
	newValue := forgedomain.ParseBitbucketAPIToken(input)
	if args.Global.Equal(newValue) {
		// the user has entered the global value --> keep using the global value, don't store the local value
		newValue = None[forgedomain.BitbucketAPIToken]()
	}
	fmt.Printf(messages.BitbucketAPITokenResult, dialogcomponents.FormattedOption(newValue, args.Global.IsSome(), exit))
	return newValue, exit, err
}

// BitbucketAPITokenHelp provides the help text for the Bitbucket API token dialog
// that matches the given Bitbucket variant.
func BitbucketAPITokenHelp(forgeType forgedomain.ForgeType) string {
	if forgeType == forgedomain.ForgeTypeBitbucketDatacenter {
		return bitbucketAPITokenHelpDataCenter
	}
	return bitbucketAPITokenHelpCloud
}
