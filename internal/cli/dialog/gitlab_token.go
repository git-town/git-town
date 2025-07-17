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
	gitLabTokenTitle = `GitLab API token`
	gitLabTokenHelp  = `
Git Town can update merge requests
and ship branches on GitLab for you.
To enable this,
please enter a GitLab API token.
More info at
https://www.git-town.com/preferences/gitlab-token.

If you leave this empty,
Git Town will not use the GitLab API.

`
)

// GitLabToken lets the user enter the GitHub API token.
func GitLabToken(oldValue Option[forgedomain.GitLabToken], inputs dialogcomponents.TestInputs) (Option[forgedomain.GitLabToken], dialogdomain.Exit, error) {
	text, exit, err := dialogcomponents.TextField(dialogcomponents.TextFieldArgs{
		DialogName:    "gitlab-token",
		ExistingValue: oldValue.String(),
		Help:          gitLabTokenHelp,
		Prompt:        "Your GitLab API token: ",
		TestInputs:    inputs,
		Title:         gitLabTokenTitle,
	})
	fmt.Printf(messages.GitLabToken, dialogcomponents.FormattedSecret(text, exit))
	return forgedomain.ParseGitLabToken(text), exit, err
}
