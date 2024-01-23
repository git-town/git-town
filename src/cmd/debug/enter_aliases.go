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
			existing := configdomain.Aliases{
				configdomain.AliasableCommandAppend: "town append",
				configdomain.AliasableCommandHack:   "town hack",
				configdomain.AliasableCommandRepo:   "other command",
			}
			dialogTestInputs := dialog.LoadTestInputs(os.Environ())
			_, _, err := dialog.Aliases(all, existing, dialogTestInputs.Next())
			return err
		},
	}
}
