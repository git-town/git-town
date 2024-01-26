package enter

import (
	"bufio"
	"fmt"
	"os"

	"github.com/git-town/git-town/v11/src/cli/dialogs/dialog"
	"github.com/git-town/git-town/v11/src/config/configdomain"
)

const enterGitLabTokenHelp = `
If you have an API token for GitLab,
and want to ship branches from the CLI,
please enter it now.
Press enter when done.
It's okay to leave this empty.

Your GitLab API token: `

// EnterGitLabToken lets the user enter the GitLab API token.
func EnterGitLabToken(inputs dialog.TestInput) (configdomain.GitLabToken, error) {
	fmt.Print(enterGitLabTokenHelp)
	reader := bufio.NewReader(os.Stdin)
	if len(inputs) > 0 {
		return configdomain.GitLabToken(inputs.ForReadline()), nil
	}
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return configdomain.GitLabToken(input), err
}
