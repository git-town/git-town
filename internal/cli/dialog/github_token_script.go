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
	githubTokenScriptTitle = `GitHub token script`
	gitHubTokenScriptHelp  = `
Enter the script that Git Town should execute
to retrieve the GitHub API token to use.

More details:
https://www.git-town.com/preferences/github-token-script

If you leave this blank,
Git Town will not interact
with the GitHub API.

`
)

func GitHubTokenScript(args Args[forgedomain.GitHubTokenScript]) (Option[forgedomain.GitHubTokenScript], dialogdomain.Exit, error) {
	input, exit, err := dialogcomponents.TextField(dialogcomponents.TextFieldArgs{
		DialogName:    "github-token",
		ExistingValue: args.Local.Or(args.Global).String(),
		Help:          gitHubTokenScriptHelp,
		Inputs:        args.Inputs,
		Prompt:        messages.GitHubTokenScriptPrompt,
		Title:         githubTokenScriptTitle,
	})
	newValue := forgedomain.ParseGitHubTokenScript(input)
	if args.Global.Equal(newValue) {
		// the user has entered the global value --> keep using the global value, don't store the local value
		newValue = None[forgedomain.GitHubTokenScript]()
	}
	fmt.Printf(messages.GitHubTokenScriptResult, dialogcomponents.FormattedSecret(newValue.String(), exit))
	return newValue, exit, err
}
