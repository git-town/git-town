package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v15/internal/cli/dialog/components"
	"github.com/git-town/git-town/v15/internal/config/configdomain"
	. "github.com/git-town/git-town/v15/internal/gohacks/prelude"
	"github.com/git-town/git-town/v15/internal/messages"
)

const (
	gitLabTokenTitle = `GitLab API token`
	gitLabTokenHelp  = `
If you have an API token for GitLab,
and want to ship branches from the CLI,
please enter it now.

It's okay to leave this empty.

`
)

// GitLabToken lets the user enter the GitHub API token.
func GitLabToken(oldValue Option[configdomain.GitLabToken], inputs components.TestInput) (Option[configdomain.GitLabToken], bool, error) {
	text, aborted, err := components.TextField(components.TextFieldArgs{
		ExistingValue: oldValue.String(),
		Help:          gitLabTokenHelp,
		Prompt:        "Your GitLab API token: ",
		TestInput:     inputs,
		Title:         gitLabTokenTitle,
	})
	fmt.Printf(messages.GitLabToken, components.FormattedSecret(text, aborted))
	return configdomain.ParseGitLabToken(text), aborted, err
}
