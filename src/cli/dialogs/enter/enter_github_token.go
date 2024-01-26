package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v11/src/config/configdomain"
)

const enterGitHubTokenHelp = `
If you have an API token for GitHub,
and want to ship branches from the CLI,
please enter it now.

Press enter when done.
It's okay to leave this empty.

`

// EnterGitHubToken lets the user enter the GitHub API token.
func EnterGitHubToken(inputs TestInput) (configdomain.GitHubToken, error) {
	fmt.Print(enterGitHubTokenHelp)
	token, _, err := textInput("existing", enterGitHubTokenHelp, "placeholder", nil)
	// reader := bufio.NewReader(os.Stdin)
	// if len(inputs) > 0 {
	// 	return configdomain.GitHubToken(inputs.ForReadline()), nil
	// }
	// input, err := reader.ReadString('\n')
	// if err != nil {
	// 	return "", err
	// }
	return configdomain.GitHubToken(token), err
}
