package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v16/internal/cli/dialog/components"
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/messages"
	. "github.com/git-town/git-town/v16/pkg/prelude"
)

const (
	githubTokenTitle = `GitHub API token`
	gitHubTokenHelp  = `
Git Town can update pull requests and ship branches on GitHub for you.
To enable this, please enter a GitHub API token.
More info at https://www.git-town.com/preferences/github-token.

If you leave this empty, Git Town will not use the GitHub API.

`
)

// GitHubToken lets the user enter the GitHub API token.
func GitHubToken(oldValue Option[configdomain.GitHubToken], inputs components.TestInput) (Option[configdomain.GitHubToken], bool, error) {
	text, aborted, err := components.TextField(components.TextFieldArgs{
		ExistingValue: oldValue.String(),
		Help:          gitHubTokenHelp,
		Prompt:        "Your GitHub API token: ",
		TestInput:     inputs,
		Title:         githubTokenTitle,
	})
	fmt.Printf(messages.GitHubToken, components.FormattedSecret(text, aborted))
	return configdomain.ParseGitHubToken(text), aborted, err
}
