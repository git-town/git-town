package dialog

import (
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/messages"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

const (
	originHostnameTitle = `Origin hostname`
	originHostnameHelp  = `
If you're using SSH identities,
specify the hostname
of your source code repository.

Only update this
if Git Town's auto-detection doesn't work.

`
)

func OriginHostname(args CommonArgs) (Option[configdomain.HostingOriginHostname], dialogdomain.Exit, error) {
	return ConfigStringDialog(ConfigStringDialogArgs[configdomain.HostingOriginHostname]{
		ConfigFileValue: args.ConfigFile.HostingOriginHostname,
		HelpText:        originHostnameHelp,
		Inputs:          args.Inputs,
		LocalValue:      args.LocalGitConfig.HostingOriginHostname,
		ParseFunc:       WrapParseFunc(configdomain.ParseHostingOriginHostname),
		PrintResultFunc: dialogcomponents.FormattedSelection,
		Prompt:          "Origin hostname: ",
		ResultMessage:   messages.OriginHostname,
		Title:           originHostnameTitle,
		UnscopedValue:   args.UnscopedGitConfig.HostingOriginHostname,
	})
}
