package dialog

import (
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/messages"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

const (
	bitbucketUsernameTitle = `Bitbucket username`
	bitbucketUsernameHelp  = `
Git Town can update pull requests
and ship branches on Bitbucket for you.
To enable this,
please enter your Bitbucket username.

If you leave this empty,
Git Town will not use the Bitbucket API.

`
)

func BitbucketUsername(args CommonArgs) (Option[forgedomain.BitbucketUsername], dialogdomain.Exit, error) {
	return ConfigStringDialog(ConfigStringDialogArgs[forgedomain.BitbucketUsername]{
		ConfigFileValue: args.ConfigFile.BitbucketUsername,
		HelpText:        bitbucketUsernameHelp,
		Inputs:          args.Inputs,
		LocalValue:      args.LocalGitConfig.BitbucketUsername,
		ParseFunc:       WrapParseFunc(forgedomain.ParseBitbucketUsername),
		PrintResultFunc: dialogcomponents.FormattedSelection,
		Prompt:          "Bitbucket username: ",
		ResultMessage:   messages.BitbucketUsername,
		Title:           bitbucketUsernameTitle,
		UnscopedValue:   args.UnscopedGitConfig.BitbucketUsername,
	})
}
