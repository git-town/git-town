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

func enterBitbucketAppPassword() *cobra.Command {
	return &cobra.Command{
		Use: "bitbucket-app-password",
		RunE: func(_ *cobra.Command, _ []string) error {
			dialogInputs := dialogcomponents.LoadTestInputs(os.Environ())
			_, _, err := dialog.ConfigStringDialog(dialog.ConfigStringDialogArgs[forgedomain.BitbucketAppPassword]{
				ConfigFileValue: None[forgedomain.BitbucketAppPassword](),
				HelpText:        dialog.BitbucketAppPasswordHelp,
				Inputs:          dialogInputs,
				LocalValue:      None[forgedomain.BitbucketAppPassword](),
				ParseFunc:       dialog.WrapParseFunc(forgedomain.ParseBitbucketAppPassword),
				Prompt:          "Your bitbucket AppPassword: ",
				ResultMessage:   messages.BitbucketAppPassword,
				Title:           dialog.BitbucketAppPasswordTitle,
				UnscopedValue:   None[forgedomain.BitbucketAppPassword](),
			})
			return err
		},
	}
}
