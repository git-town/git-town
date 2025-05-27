package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/cli/dialog/components"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/messages"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

const (
	githubTokenTitle = `GitHub API token`
	gitHubTokenHelp  = `
Git Town can update pull requests
and ship branches on your behalf
using the GitHub API.
To enable this, enter a GitHub API token.

More details: https://www.git-town.com/preferences/github-token

If you leave this blank,
Git Town will not interact with the GitHub API.

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
