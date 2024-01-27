package enter

import (
	"fmt"

	"github.com/git-town/git-town/v11/src/cli/dialogs/dialog"
	"github.com/git-town/git-town/v11/src/config/configdomain"
)

const enterGitLabTokenHelp = `
If you have an API token for GitLab,
and want to ship branches from the CLI,
please enter it now.

It's okay to leave this empty.

`

// GitLabToken lets the user enter the GitHub API token.
func GitLabToken(oldValue configdomain.GitLabToken, inputs dialog.TestInput) (configdomain.GitLabToken, bool, error) {
	token, aborted, err := dialog.TextField(oldValue.String(), enterGitLabTokenHelp, "Your GitLab API token: ", inputs)
	fmt.Printf("GitLab token: %s\n", dialog.FormattedToken(token, aborted))
	return configdomain.GitLabToken(token), aborted, err
}
