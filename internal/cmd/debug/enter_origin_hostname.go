package debug

import (
	"os"

	"github.com/git-town/git-town/v21/internal/cli/dialog"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/messages"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/spf13/cobra"
)

func enterOriginHostname() *cobra.Command {
	return &cobra.Command{
		Use: "origin-hostname",
		RunE: func(_ *cobra.Command, _ []string) error {
			dialogInputs := dialogcomponents.LoadTestInputs(os.Environ())
			_, _, err := dialog.ConfigStringDialog(dialog.ConfigStringDialogArgs[configdomain.HostingOriginHostname]{
				ConfigFileValue: None[configdomain.HostingOriginHostname](),
				HelpText:        dialog.OriginHostnameHelp,
				Inputs:          dialogInputs,
				LocalValue:      None[configdomain.HostingOriginHostname](),
				ParseFunc:       dialog.WrapParseFunc(configdomain.ParseHostingOriginHostname),
				Prompt:          "Your origin hostname: ",
				ResultMessage:   messages.OriginHostname,
				Title:           dialog.OriginHostnameTitle,
				UnscopedValue:   None[configdomain.HostingOriginHostname](),
			})
			return err
		},
	}
}
