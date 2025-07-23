package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/messages"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

const (
	githubTokenTitle = `GitHub API token`
	gitHubTokenHelp  = `
Git Town can update pull requests
and ship branches on your behalf
using the GitHub API.
To enable this,
enter a GitHub API token.

More details:
https://www.git-town.com/preferences/github-token

If you leave this blank,
Git Town will not interact
with the GitHub API.

`
)

func GitHubToken(args Args[forgedomain.GitHubToken]) (Option[forgedomain.GitHubToken], dialogdomain.Exit, error) {
	input, exit, err := dialogcomponents.TextField(dialogcomponents.TextFieldArgs{
		DialogName:    "github-token",
		ExistingValue: args.Local.Or(args.Global).String(),
		Help:          gitHubTokenHelp,
		Inputs:        args.Inputs,
		Prompt:        messages.GitHubTokenPrompt,
		Title:         githubTokenTitle,
	})
	newValue := forgedomain.ParseGitHubToken(input)
	if args.Global.Equal(newValue) {
		// the user has entered the global value --> keep using the global value, don't store the local value
		newValue = None[forgedomain.GitHubToken]()
	}
	fmt.Printf(messages.GitHubTokenResult, dialogcomponents.FormattedSecret(newValue.String(), exit))
	return newValue, exit, err
}
