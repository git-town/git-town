package dialog

import (
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/messages"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

const (
	codebergTokenTitle = `Codeberg API token`
	codebergTokenHelp  = `
Git Town can update pull requests
and ship branches on codeberg-based forges for you.
To enable this, please enter a codeberg API token.
More info at
https://docs.codeberg.org/advanced/access-token.

If you leave this empty,
Git Town will not use the codeberg API.

`
)

func CodebergToken(args CommonArgs) (Option[forgedomain.CodebergToken], dialogdomain.Exit, error) {
	return ConfigStringDialog(ConfigStringDialogArgs[forgedomain.CodebergToken]{
		ConfigFileValue: args.ConfigFile.CodebergToken,
		HelpText:        codebergTokenHelp,
		Inputs:          args.Inputs,
		LocalValue:      args.LocalGitConfig.CodebergToken,
		ParseFunc:       WrapParseFunc(forgedomain.ParseCodebergToken),
		PrintResultFunc: dialogcomponents.FormattedSecret,
		Prompt:          "Codeberg token: ",
		ResultMessage:   messages.CodebergToken,
		Title:           codebergTokenTitle,
		UnscopedValue:   args.UnscopedGitConfig.CodebergToken,
	})
}
