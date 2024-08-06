package debug

import (
	"os"

	"github.com/git-town/git-town/v14/internal/cli/dialog"
	"github.com/git-town/git-town/v14/internal/cli/dialog/components"
	"github.com/git-town/git-town/v14/internal/config/configdomain"
	"github.com/git-town/git-town/v14/pkg/keys"
	"github.com/spf13/cobra"
)

func enterAliases() *cobra.Command {
	return &cobra.Command{
		Use: "aliases",
		RunE: func(_ *cobra.Command, _ []string) error {
			all := keys.AllAliasableCommands()
			existing := configdomain.Aliases{
				keys.AliasableCommandAppend: "town append",
				keys.AliasableCommandHack:   "town hack",
				keys.AliasableCommandRepo:   "other command",
			}
			dialogTestInputs := components.LoadTestInputs(os.Environ())
			_, _, err := dialog.Aliases(all, existing, dialogTestInputs.Next())
			return err
		},
	}
}
