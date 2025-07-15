package dialog

import (
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/messages"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

const (
	giteaTokenTitle = `Gitea API token`
	giteaTokenHelp  = `
Git Town can update pull requests
and ship branches on gitea for you.
To enable this, please enter a gitea API token.
More info at
https://www.git-town.com/preferences/gitea-token.

If you leave this empty,
Git Town will not use the gitea API.

`
)

func GiteaToken(args CommonArgs) (Option[forgedomain.GiteaToken], dialogdomain.Exit, error) {
	return ConfigStringDialog(ConfigStringDialogArgs[forgedomain.GiteaToken]{
		ConfigFileValue: args.ConfigFile.GiteaToken,
		HelpText:        giteaTokenHelp,
		Inputs:          args.Inputs,
		LocalValue:      args.LocalGitConfig.GiteaToken,
		ParseFunc:       WrapParseFunc(forgedomain.ParseGiteaToken),
		Prompt:          "Gitea token: ",
		ResultMessage:   messages.GiteaToken,
		Title:           giteaTokenTitle,
		UnscopedValue:   args.UnscopedGitConfig.GiteaToken,
	})
}
