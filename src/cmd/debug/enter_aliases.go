package debug

import (
	"os"

	"github.com/git-town/git-town/v12/src/cli/dialog"
	"github.com/git-town/git-town/v12/src/cli/dialog/components"
	"github.com/git-town/git-town/v12/src/config/configdomain"
	"github.com/spf13/cobra"
)

func enterAliases() *cobra.Command {
	return &cobra.Command{
		Use: "aliases",
		RunE: func(_ *cobra.Command, _ []string) error {
			all := configdomain.AllAliasableCommands()
			existing := configdomain.Aliases{
				configdomain.AliasableCommandAppend: "town append",
				configdomain.AliasableCommandHack:   "town hack",
				configdomain.AliasableCommandRepo:   "other command",
			}
			dialogTestInputs := components.LoadTestInputs(os.Environ())
			_, _, err := dialog.Aliases(all, existing, dialogTestInputs.Next())
			return err
		},
	}
}
