package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v16/internal/cli/dialog/components"
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/messages"
	. "github.com/git-town/git-town/v16/pkg/prelude"
)

//nolint:gosec
const (
	bitbucketAppPasswordTitle = `Bitbucket App Password/Token`
	bitbucketAppPasswordHelp  = `
Git Town can update pull requests and ship branches on Bitbucket for you.
To enable this, please enter a Bitbucket App Password or token.
This is not your normal account password.
More info at https://www.git-town.com/preferences/bitbucket-app-password.

If you leave this empty, Git Town will not use the Bitbucket API.

`
)

// GitHubToken lets the user enter the GitHub API token.
func BitbucketAppPassword(oldValue Option[configdomain.BitbucketAppPassword], inputs components.TestInput) (Option[configdomain.BitbucketAppPassword], bool, error) {
	text, aborted, err := components.TextField(components.TextFieldArgs{
		ExistingValue: oldValue.String(),
		Help:          bitbucketAppPasswordHelp,
		Prompt:        "Bitbucket App Password/Token: ",
		TestInput:     inputs,
		Title:         bitbucketAppPasswordTitle,
	})
	fmt.Printf(messages.GitHubToken, components.FormattedSecret(text, aborted))
	return configdomain.ParseBitbucketAppPassword(text), aborted, err
}
