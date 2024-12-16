package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v17/internal/cli/dialog/components"
	"github.com/git-town/git-town/v17/internal/config/configdomain"
	"github.com/git-town/git-town/v17/internal/messages"
	. "github.com/git-town/git-town/v17/pkg/prelude"
)

const (
	gitLabTokenTitle = `GitLab API token`
	gitLabTokenHelp  = `
Git Town can update merge requests and ship branches on GitLab for you.
To enable this, please enter a GitLab API token.
More info at https://www.git-town.com/preferences/gitlab-token.

If you leave this empty, Git Town will not use the GitLab API.

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
