package dialog

import (
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/messages"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

const (
	gitHubTokenTitle = `GitHub API token`
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

func GitHubToken(args CommonArgs) (Option[forgedomain.GitHubToken], dialogdomain.Exit, error) {
	return ConfigStringDialog(ConfigStringDialogArgs[forgedomain.GitHubToken]{
		ConfigFileValue: args.ConfigFile.GitHubToken,
		HelpText:        gitHubTokenHelp,
		Inputs:          args.Inputs,
		LocalValue:      args.LocalGitConfig.GitHubToken,
		ParseFunc:       WrapParseFunc(forgedomain.ParseGitHubToken),
		Prompt:          "GitHub token: ",
		ResultMessage:   messages.GitHubToken,
		Title:           gitHubTokenTitle,
		UnscopedValue:   args.UnscopedGitConfig.GitHubToken,
	})
}
