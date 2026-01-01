package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/messages"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

const (
	githubTokenTitle = `GitHub API token`
	githubTokenHelp  = `
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

func GithubToken(args Args[forgedomain.GithubToken]) (Option[forgedomain.GithubToken], dialogdomain.Exit, error) {
	input, exit, err := dialogcomponents.TextField(dialogcomponents.TextFieldArgs{
		DialogName:    "github-token",
		ExistingValue: args.Local.Or(args.Global).StringOr(""),
		Help:          githubTokenHelp,
		Inputs:        args.Inputs,
		Prompt:        messages.GithubTokenPrompt,
		Title:         githubTokenTitle,
	})
	newValue := forgedomain.ParseGitHubToken(input)
	if args.Global.Equal(newValue) {
		// the user has entered the global value --> keep using the global value, don't store the local value
		newValue = None[forgedomain.GithubToken]()
	}
	fmt.Printf(messages.GithubTokenResult, dialogcomponents.FormattedSecret(newValue.StringOr(""), exit))
	return newValue, exit, err
}
