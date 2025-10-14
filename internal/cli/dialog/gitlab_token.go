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

func GitLabToken(args Args[forgedomain.GitLabToken]) (Option[forgedomain.GitLabToken], dialogdomain.Exit, error) {
	input, exit, err := dialogcomponents.TextField(dialogcomponents.TextFieldArgs{
		DialogName:    "gitlab-token",
		ExistingValue: args.Local.Or(args.Global).StringOr(""),
		Help:          gitLabTokenHelp,
		Inputs:        args.Inputs,
		Prompt:        messages.GitLabTokenPrompt,
		Title:         gitLabTokenTitle,
	})
	newValue := forgedomain.ParseGitLabToken(input)
	if args.Global.Equal(newValue) {
		// the user has entered the global value --> keep using the global value, don't store the local value
		newValue = None[forgedomain.GitLabToken]()
	}
	fmt.Printf(messages.GitLabTokenResult, dialogcomponents.FormattedSecret(newValue.GetOrZero().String(), exit))
	return newValue, exit, err
}
