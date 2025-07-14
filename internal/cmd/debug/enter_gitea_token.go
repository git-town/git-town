package debug

import (
	"os"

	"github.com/git-town/git-town/v21/internal/cli/dialog"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/messages"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/spf13/cobra"
)

func enterGiteaToken() *cobra.Command {
	return &cobra.Command{
		Use: "gitea-token",
		RunE: func(_ *cobra.Command, _ []string) error {
			dialogInputs := dialogcomponents.LoadTestInputs(os.Environ())
			_, _, err := dialog.ConfigStringDialog(dialog.ConfigStringDialogArgs[forgedomain.GiteaToken]{
				ConfigFileValue: None[forgedomain.GiteaToken](),
				HelpText:        dialog.GiteaTokenHelp,
				Inputs:          dialogInputs,
				LocalValue:      None[forgedomain.GiteaToken](),
				ParseFunc:       dialog.WrapParseFunc(forgedomain.ParseGiteaToken),
				Prompt:          "Your Gitea token: ",
				ResultMessage:   messages.GiteaToken,
				Title:           dialog.GiteaTokenTitle,
				UnscopedValue:   None[forgedomain.GiteaToken](),
			})
			return err
		},
	}
}
