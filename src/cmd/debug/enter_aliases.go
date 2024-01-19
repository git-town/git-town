package debug

import (
	"os"

	"github.com/git-town/git-town/v11/src/cli/dialog"
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/spf13/cobra"
)

func enterAliases() *cobra.Command {
	return &cobra.Command{
		Use: "aliases",
		RunE: func(cmd *cobra.Command, args []string) error {
			all := configdomain.AllAliasableCommands()
			existing := configdomain.AliasableCommands{
				configdomain.AliasableCommandAppend,
				configdomain.AliasableCommandHack,
			}
			dialogTestInputs := dialog.LoadTestInputs(os.Environ())
			_, _, err := dialog.EnterAliases(all, existing, dialogTestInputs.Next())
			return err
		},
	}
}
