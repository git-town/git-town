package dialogs

import (
	"fmt"

	"github.com/git-town/git-town/v11/src/cli/dialogs/components"
	"github.com/git-town/git-town/v11/src/config/configdomain"
)

const enterGitHubTokenHelp = `
If you have an API token for GitHub,
and want to ship branches from the CLI,
please enter it now.

It's okay to leave this empty.

`

// GitHubToken lets the user enter the GitHub API token.
func GitHubToken(oldValue configdomain.GitHubToken, inputs components.TestInput) (configdomain.GitHubToken, bool, error) {
	token, aborted, err := components.TextField(components.TextFieldArgs{
		ExistingValue: oldValue.String(),
		Help:          enterGitHubTokenHelp,
		Prompt:        "Your GitHub API token: ",
		TestInput:     inputs,
	})
	fmt.Printf("GitHub token: %s\n", components.FormattedToken(token, aborted))
	return configdomain.GitHubToken(token), aborted, err
}
