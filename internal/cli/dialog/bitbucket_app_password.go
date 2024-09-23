package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v16/internal/cli/dialog/components"
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/messages"
	. "github.com/git-town/git-town/v16/pkg/prelude"
)

const (
	bitbucketAppPasswordTitle = `Bitbucket App Password`
	bitbucketAppPasswordHelp  = `
If you want Git Town to use the Bitbucket API,
please enter your Bitbucket App Password.
More info at https://support.atlassian.com/bitbucket-cloud/docs/app-passwords.

It's okay to leave this empty.

`
)

// GitHubToken lets the user enter the GitHub API token.
func BitbucketAppPassword(oldValue Option[configdomain.BitbucketAppPassword], inputs components.TestInput) (Option[configdomain.BitbucketAppPassword], bool, error) {
	text, aborted, err := components.TextField(components.TextFieldArgs{
		ExistingValue: oldValue.String(),
		Help:          bitbucketAppPasswordHelp,
		Prompt:        "Your Bitbucket App Password: ",
		TestInput:     inputs,
		Title:         bitbucketAppPasswordTitle,
	})
	fmt.Printf(messages.GitHubToken, components.FormattedSecret(text, aborted))
	return configdomain.ParseBitbucketAppPassword(text), aborted, err
}
