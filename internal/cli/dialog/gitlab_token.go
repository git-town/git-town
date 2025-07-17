package dialog

import (
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

func GitLabToken(args CommonArgs) (Option[forgedomain.GitLabToken], dialogdomain.Exit, error) {
	return ConfigStringDialog(ConfigStringDialogArgs[forgedomain.GitLabToken]{
		ConfigFileValue: args.ConfigFile.GitLabToken,
		DialogName:      "gitlab-token",
		HelpText:        gitLabTokenHelp,
		Inputs:          args.Inputs,
		LocalValue:      args.LocalGitConfig.GitLabToken,
		ParseFunc:       WrapParseFunc(forgedomain.ParseGitLabToken),
		PrintResultFunc: dialogcomponents.FormattedSecret,
		Prompt:          "GitLab token: ",
		ResultMessage:   messages.GitLabToken,
		Title:           gitLabConnectorTypeTitle,
		UnscopedValue:   args.UnscopedGitConfig.GitLabToken,
	})
}
