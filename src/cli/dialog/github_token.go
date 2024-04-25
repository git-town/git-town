package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/gohacks"
	"github.com/git-town/git-town/v14/src/messages"
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
func GitHubToken(oldValue gohacks.Option[configdomain.GitHubToken], inputs components.TestInput) (*configdomain.GitHubToken, bool, error) {
	var tokenText string
	if value, has := oldValue.Get(); has {
		tokenText = value.String()
	}
	text, aborted, err := components.TextField(components.TextFieldArgs{
		ExistingValue: tokenText,
		Help:          gitHubTokenHelp,
		Prompt:        "Your GitHub API token: ",
		TestInput:     inputs,
		Title:         githubTokenTitle,
	})
	fmt.Printf(messages.GitHubToken, components.FormattedSecret(text, aborted))
	return configdomain.NewGitHubTokenRef(text), aborted, err
}
