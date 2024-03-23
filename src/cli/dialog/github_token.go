package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v13/src/cli/dialog/components"
	"github.com/git-town/git-town/v13/src/config/configdomain"
	"github.com/git-town/git-town/v13/src/messages"
)

const (
	githubTokenTitle = `GitHub API token`
	gitHubTokenHelp  = `
If you have an API token for GitHub,
and want to ship branches from the CLI,
please enter it now.

It's okay to leave this empty.

`
)

// GitHubToken lets the user enter the GitHub API token.
func GitHubToken(oldValue configdomain.GitHubToken, inputs components.TestInput) (configdomain.GitHubToken, bool, error) {
	token, aborted, err := components.TextField(components.TextFieldArgs{
		ExistingValue: oldValue.String(),
		Help:          gitHubTokenHelp,
		Prompt:        "Your GitHub API token: ",
		TestInput:     inputs,
		Title:         githubTokenTitle,
	})
	fmt.Printf(messages.GitHubToken, components.FormattedSecret(token, aborted))
	return configdomain.GitHubToken(token), aborted, err
}
