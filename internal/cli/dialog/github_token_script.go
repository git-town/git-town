package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/cli/dialog/components"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/messages"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

const (
	githubTokenScriptTitle = `GitHub API token script`
	gitHubTokenScriptHelp  = `
Git Town can update pull requests
and ship branches on your behalf
using the GitHub API.
To enable this,
enter a script that provides
the GitHub API token.

More details:
https://www.git-town.com/preferences/github-token-script

Please enter the script that provides the GitHub API token below.

`
)

// GitHubTokenScript lets the user enter the GitHub API token.
func GitHubTokenScript(oldValue Option[forgedomain.GitHubTokenScript], inputs components.TestInput) (Option[forgedomain.GitHubTokenScript], dialogdomain.Exit, error) {
	text, exit, err := components.TextField(components.TextFieldArgs{
		ExistingValue: oldValue.String(),
		Help:          gitHubTokenScriptHelp,
		Prompt:        "GitHub token script: ",
		TestInput:     inputs,
		Title:         githubTokenScriptTitle,
	})
	fmt.Printf(messages.GitHubTokenScript, components.FormattedSecret(text, exit))
	return NewOption(forgedomain.GitHubTokenScript(text)), exit, err

}
