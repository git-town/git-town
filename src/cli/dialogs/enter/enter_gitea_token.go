package dialog

import (
	"bufio"
	"fmt"
	"os"

	"github.com/git-town/git-town/v11/src/config/configdomain"
)

const enterGiteaTokenHelp = `
If you have an API token for Gitea,
and want to ship branches from the CLI,
please enter it now.
Press enter when done.
It's okay to leave this empty.

Your Gitea API token: `

// EnterGiteaToken lets the user enter the Gitea API token.
func EnterGiteaToken(inputs TestInput) (configdomain.GiteaToken, error) {
	fmt.Print(enterGiteaTokenHelp)
	reader := bufio.NewReader(os.Stdin)
	if len(inputs) > 0 {
		return configdomain.GiteaToken(inputs.ForReadline()), nil
	}
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return configdomain.GiteaToken(input), err
}
