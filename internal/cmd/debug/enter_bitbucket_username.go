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

func enterBitbucketUsername() *cobra.Command {
	return &cobra.Command{
		Use: "bitbucket-username",
		RunE: func(_ *cobra.Command, _ []string) error {
			dialogInputs := dialogcomponents.LoadTestInputs(os.Environ())
			_, _, err := dialog.ConfigStringDialog(dialog.ConfigStringDialogArgs[forgedomain.BitbucketUsername]{
				ConfigFileValue: None[forgedomain.BitbucketUsername](),
				HelpText:        dialog.BitbucketUsernameHelp,
				Inputs:          dialogInputs,
				LocalValue:      None[forgedomain.BitbucketUsername](),
				ParseFunc:       dialog.WrapParseFunc(forgedomain.ParseBitbucketUsername),
				Prompt:          "Your bitbucket username: ",
				ResultMessage:   messages.BitbucketUsername,
				Title:           dialog.BitbucketUsernameTitle,
				UnscopedValue:   None[forgedomain.BitbucketUsername](),
			})
			return err
		},
	}
}
