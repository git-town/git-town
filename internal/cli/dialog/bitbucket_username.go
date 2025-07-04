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

func BitbucketUsername(oldValue Option[forgedomain.BitbucketUsername], inputs dialogcomponents.TestInput) (Option[forgedomain.BitbucketUsername], dialogdomain.Exit, error) {
	text, exit, err := dialogcomponents.TextField(dialogcomponents.TextFieldArgs{
		ExistingValue: oldValue.String(),
		Help:          bitbucketUsernameHelp,
		Prompt:        "Your Bitbucket username: ",
		TestInput:     inputs,
		Title:         bitbucketUsernameTitle,
	})
	fmt.Printf(messages.BitbucketUsername, dialogcomponents.FormattedSecret(text, exit))
	return forgedomain.ParseBitbucketUsername(text), exit, err
}
