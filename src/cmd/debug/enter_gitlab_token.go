package debug

import (
	"os"

	"github.com/git-town/git-town/v11/src/cli/dialogs/dialog"
	"github.com/git-town/git-town/v11/src/cli/dialogs/enter"
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/spf13/cobra"
)

func enterGitLabToken() *cobra.Command {
	return &cobra.Command{
		Use: "gitlab-token",
		RunE: func(cmd *cobra.Command, args []string) error {
			dialogInputs := dialog.LoadTestInputs(os.Environ())
			_, _, err := enter.GitLabToken(configdomain.GitLabToken(""), dialogInputs.Next())
			return err
		},
	}
}
