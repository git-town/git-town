package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v18/internal/cli/dialog/components"
	"github.com/git-town/git-town/v18/internal/config/configdomain"
	"github.com/git-town/git-town/v18/internal/messages"
	. "github.com/git-town/git-town/v18/pkg/prelude"
)

const (
	bitbucketUsernameTitle = `Bitbucket username`
	bitbucketUsernameHelp  = `
Git Town can update pull requests and ship branches on Bitbucket for you.
To enable this, please enter your Bitbucket username.

If you leave this empty, Git Town will not use the Bitbucket API.

`
)

// BitbucketUsername lets the user enter the Bitbucket username.
func BitbucketUsername(oldValue Option[configdomain.BitbucketUsername], inputs components.TestInput) (Option[configdomain.BitbucketUsername], bool, error) {
	text, aborted, err := components.TextField(components.TextFieldArgs{
		ExistingValue: oldValue.String(),
		Help:          bitbucketUsernameHelp,
		Prompt:        "Your Bitbucket username: ",
		TestInput:     inputs,
		Title:         bitbucketUsernameTitle,
	})
	fmt.Printf(messages.BitbucketUsername, components.FormattedSecret(text, aborted))
	return configdomain.ParseBitbucketUsername(text), aborted, err
}
