package enter

import (
	"fmt"

	"github.com/git-town/git-town/v11/src/cli/dialogs/dialog"
	"github.com/git-town/git-town/v11/src/config/configdomain"
)

const enterGitHubTokenHelp = `
If you have an API token for GitHub,
and want to ship branches from the CLI,
please enter it now.

It's okay to leave this empty.

`

// GitHubToken lets the user enter the GitHub API token.
func GitHubToken(oldValue configdomain.GitHubToken, inputs dialog.TestInput) (configdomain.GitHubToken, bool, error) {
	token, aborted, err := dialog.TextField(oldValue.String(), enterGitHubTokenHelp, "Your GitHub API token: ", inputs)
	fmt.Printf("GitHub token: %s\n", dialog.FormattedToken(token, aborted))
	return configdomain.GitHubToken(token), aborted, err
}
